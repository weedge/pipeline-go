package processors

import (
	"testing"
)

func TestObjID(t *testing.T) {
	// Reset the global ID counter for testing
	_ID = 0

	id1 := ObjID()
	id2 := ObjID()
	id3 := ObjID()

	if id1 != 1 {
		t.Errorf("Expected first ID to be 1, got %d", id1)
	}

	if id2 != 2 {
		t.Errorf("Expected second ID to be 2, got %d", id2)
	}

	if id3 != 3 {
		t.Errorf("Expected third ID to be 3, got %d", id3)
	}
}

func TestObjCount(t *testing.T) {
	// Reset the counts map for testing
	_COUNTS = make(map[string]int)

	type TestStruct struct{}
	type AnotherStruct struct{}

	obj1 := &TestStruct{}
	obj2 := &TestStruct{}
	obj3 := &AnotherStruct{}

	// Test counting for TestStruct
	count1 := ObjCount(obj1)
	if count1 != 1 {
		t.Errorf("Expected first count of TestStruct to be 1, got %d", count1)
	}

	count2 := ObjCount(obj2)
	if count2 != 2 {
		t.Errorf("Expected second count of TestStruct to be 2, got %d", count2)
	}

	// Test counting for AnotherStruct
	count3 := ObjCount(obj3)
	if count3 != 1 {
		t.Errorf("Expected first count of AnotherStruct to be 1, got %d", count3)
	}

	// Test that TestStruct count is still 2
	count4 := ObjCount(obj1)
	if count4 != 3 {
		t.Errorf("Expected third count of TestStruct to be 3, got %d", count4)
	}
}

func TestGetTypeName(t *testing.T) {
	type TestStruct struct{}

	var nilPtr *TestStruct
	obj := &TestStruct{}

	// Test with a regular object
	name := getTypeName(obj)
	if name != "TestStruct" {
		t.Errorf("Expected type name to be 'TestStruct', got '%s'", name)
	}

	// Test with a nil pointer (this would panic without our fix)
	name = getTypeName(nilPtr)
	if name != "TestStruct" {
		t.Errorf("Expected type name to be 'TestStruct', got '%s'", name)
	}
}