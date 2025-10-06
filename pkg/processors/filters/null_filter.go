package filters

import (
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// NullFilter drops all frames that pass through it, except for control frames.
type NullFilter struct {
	processors.FrameProcessor
}

// NewNullFilter creates a new NullFilter.
func NewNullFilter() *NullFilter {
	return &NullFilter{}
}

// ProcessFrame drops all non-control frames.
func (f *NullFilter) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Only pass control frames through.
	switch frame.(type) {
	case *frames.StartFrame, *frames.EndFrame, *frames.CancelFrame, *frames.StopTaskFrame, *frames.MetricsFrame, *frames.SyncFrame:
		f.PushFrame(frame, direction)
	case *frames.ErrorFrame: // ErrorFrame has pointer receiver for String method
		f.PushFrame(frame, direction)
	default:
		// Drop all other frames.
	}
}
