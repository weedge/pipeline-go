package filters

import (
	"reflect"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// TypeFilter drops frames based on their type.
type TypeFilter struct {
	processors.FrameProcessor
	includeTypes []reflect.Type
}

// NewTypeFilter creates a new TypeFilter.
func NewTypeFilter(includeTypes []interface{}) *TypeFilter {
	var types []reflect.Type
	for _, t := range includeTypes {
		types = append(types, reflect.TypeOf(t))
	}
	return &TypeFilter{
		includeTypes: types,
	}
}

// isControlFrame checks if a frame is a control frame that should always pass through.
func (f *TypeFilter) isControlFrame(frame frames.Frame) bool {
	switch frame.(type) {
	case *frames.StartFrame, *frames.EndFrame, *frames.CancelFrame, *frames.ErrorFrame, *frames.StopTaskFrame, *frames.MetricsFrame, *frames.SyncFrame:
		return true
	default:
		return false
	}
}

// ProcessFrame applies the filter and pushes the frame if it passes.
func (f *TypeFilter) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	if f.isControlFrame(frame) {
		f.PushFrame(frame, direction)
		return
	}

	frameType := reflect.TypeOf(frame)
	for _, t := range f.includeTypes {
		if frameType == t {
			f.PushFrame(frame, direction)
			return
		}
	}
}
