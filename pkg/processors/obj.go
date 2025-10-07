package processors

import (
	"reflect"
	"sync"
)

var (
	// _COUNTS is a map to store object counts
	_COUNTS = make(map[string]int)
	// _COUNTS_MUTEX is a mutex to protect _COUNTS
	_COUNTS_MUTEX sync.Mutex

	// _ID is a global counter for object IDs
	_ID = 0
	// _ID_MUTEX is a mutex to protect _ID
	_ID_MUTEX sync.Mutex
)

// ObjID generates and returns a unique object ID
func ObjID() int {
	_ID_MUTEX.Lock()
	defer _ID_MUTEX.Unlock()
	_ID++
	return _ID
}

// ObjCount increments and returns the count for the given object type
func ObjCount(obj interface{}) int {
	// Get the name of the object's type
	name := getTypeName(obj)

	_COUNTS_MUTEX.Lock()
	defer _COUNTS_MUTEX.Unlock()

	_COUNTS[name]++
	return _COUNTS[name]
}

// getTypeName returns the name of the type of the given object
func getTypeName(obj interface{}) string {
	// Using reflection to get the type name
	// Note: In Go, we could also use fmt.Sprintf("%T", obj) for the full type name
	t := reflect.TypeOf(obj)
	if t == nil {
		return "<nil>"
	}
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
