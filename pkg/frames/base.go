package frames

import (
	"fmt"

	"github.com/weedge/pipeline-go/pkg"
)

// Frame is the interface that all frame types implement.
type Frame interface {
	ID() uint64
	Name() string
	String() string
}

// BaseFrame is a struct that provides a default implementation of the Frame interface.
type BaseFrame struct {
	id   uint64
	name string
}

// NewBaseFrameWithName creates a new BaseFrame.
func NewBaseFrameWithName(name string) *BaseFrame {
	if name == "" {
		name = "BaseFrame"
	}
	id := pkg.CountForType(name)
	return &BaseFrame{
		id:   id,
		name: fmt.Sprintf("%s#%d", name, id),
	}
}

func NewBaseFrame() *BaseFrame {
	name := "BaseFrame"
	id := pkg.CountForType(name)
	return &BaseFrame{
		id:   id,
		name: fmt.Sprintf("%s#%d", name, id),
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

// String returns the frame's name.
func (f *BaseFrame) String() string {
	return f.name
}
