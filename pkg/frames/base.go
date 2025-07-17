package frames

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// - `Frame` Struct (Go):**
//   - I'll create a `Frame` struct with `ID` (int) and `Name` (string) fields.
//   - The `obj_id` and `obj_count` logic seems specific to the Python implementation's runtime object tracking. For the Go version, I can replicate this with a package-level counter or a simpler mechanism. A global atomic counter for IDs and a type-based counter for names would be a good Go-idiomatic equivalent.
//   - I'll create a constructor function like `NewFrame(typeName string)` that initializes the `ID` and `Name`.
var (
	globalFrameID int64
	typeCounts    = make(map[string]*int64)
	typeCountsMtx sync.Mutex
)

func nextFrameID() int64 {
	return atomic.AddInt64(&globalFrameID, 1)
}

func countForType(typeName string) int64 {
	typeCountsMtx.Lock()
	defer typeCountsMtx.Unlock()
	count, ok := typeCounts[typeName]
	if !ok {
		count = new(int64)
		typeCounts[typeName] = count
	}
	return atomic.AddInt64(count, 1)
}

var (
	nextID uint64
)

// Frame is the interface that all frame types implement.
type Frame interface {
	ID() uint64
	Name() string
}

// BaseFrame is a struct that provides a default implementation of the Frame interface.
type BaseFrame struct {
	id   uint64
	name string
}

// NewBaseFrame creates a new BaseFrame.
func NewBaseFrame(name string) *BaseFrame {
	return &BaseFrame{
		id:   atomic.AddUint64(&nextID, 1),
		name: fmt.Sprintf("%s#%d", name, nextID),
	}
}

// ID returns the frame's ID.
func (f *BaseFrame) ID() uint64 {
	return f.id
}

// Name returns the frame's name.
func (f *BaseFrame) Name() string {
	return f.name
}
