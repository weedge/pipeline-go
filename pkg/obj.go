package pkg

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	// _COUNTS is a map to store object counts
	_COUNTS = make(map[string]int)
	// _COUNTS_MUTEX is a mutex to protect _COUNTS
	_COUNTS_MUTEX sync.Mutex

	// _ID is a global counter for object IDs
	_ID uint64 = 0
)

// ObjID generates and returns a unique object ID
func ObjID() uint64 {
	return atomic.AddUint64(&_ID, 1)
}

// ObjCount increments and returns the count for the given object type
func ObjCount(obj any) int {
	// Get the name of the object's type
	name := GetObjTypeName(obj)

	_COUNTS_MUTEX.Lock()
	defer _COUNTS_MUTEX.Unlock()

	_COUNTS[name]++
	return _COUNTS[name]
}

// GetObjTypeName returns the name of the type of the given object
func GetObjTypeName(obj any) string {
	// Using reflection to get the type name
	// Note: In Go, we could also use fmt.Sprintf("%T", obj) for the full type name
	t := reflect.TypeOf(obj)
	if t == nil {
		return "<nil>"
	}
	if t.Kind() == reflect.Pointer {
		return t.Elem().Name()
	}
	return t.Name()
}

// - `Frame` Struct (Go):**
//   - I'll create a `Frame` struct with `ID` (int) and `Name` (string) fields.
//   - The `obj_id` and `obj_count` logic seems specific to the Python implementation's runtime object tracking. For the Go version, I can replicate this with a package-level counter or a simpler mechanism. A global atomic counter for IDs and a type-based counter for names would be a good Go-idiomatic equivalent.
//   - I'll create a constructor function like `NewFrame(typeName string)` that initializes the `ID` and `Name`.
var (
	typeCounts    = make(map[string]*uint64)
	typeCountsMtx sync.Mutex
)

func CountForType(typeName string) uint64 {
	typeCountsMtx.Lock()
	defer typeCountsMtx.Unlock()
	count, ok := typeCounts[typeName]
	if !ok {
		count = new(uint64)
		typeCounts[typeName] = count
	}
	return atomic.AddUint64(count, 1)
}
