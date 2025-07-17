package filters

import (
	"reflect"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// FilterMode defines whether the filter allows or denies the specified types.
type FilterMode int

const (
	// Allow mode: Only frames of the specified types will pass.
	Allow FilterMode = iota
	// Deny mode: Frames of the specified types will be dropped.
	Deny
)

// TypeFilter filters frames based on their type.
type TypeFilter struct {
	processors.BaseProcessor
	types map[reflect.Type]bool
	mode  FilterMode
}

// NewTypeFilter creates a new TypeFilter.
// The types argument is a list of frame types to filter (e.g., reflect.TypeOf(&frames.TextFrame{})).
func NewTypeFilter(mode FilterMode, types ...reflect.Type) *TypeFilter {
	typeMap := make(map[reflect.Type]bool)
	for _, t := range types {
		typeMap[t] = true
	}
	return &TypeFilter{
		types: typeMap,
		mode:  mode,
	}
}

// isControlFrame checks if a frame is a control frame that should always pass through.
func (f *TypeFilter) isControlFrame(frame frames.Frame) bool {
	switch frame.(type) {
	case frames.StartFrame, frames.EndFrame, frames.CancelFrame, frames.ErrorFrame, frames.StopTaskFrame, frames.MetricsFrame, frames.SyncFrame:
		return true
	default:
		return false
	}
}

func (f *TypeFilter) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	if f.isControlFrame(frame) {
		f.PushFrame(frame, direction)
		return
	}

	frameType := reflect.TypeOf(frame)
	_, found := f.types[frameType]

	shouldPass := false
	if f.mode == Allow {
		if found {
			shouldPass = true
		}
	} else { // Deny mode
		if !found {
			shouldPass = true
		}
	}

	if shouldPass {
		f.PushFrame(frame, direction)
	}
}
