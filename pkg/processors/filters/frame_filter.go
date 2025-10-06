package filters

import (
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// FrameFilter drops frames based on a filter function.
type FrameFilter struct {
	processors.FrameProcessor
	filterFn func(frame frames.Frame) bool
}

// NewFrameFilter creates a new FrameFilter.
// The filter function should return true to keep the frame, false to drop it.
func NewFrameFilter(filterFn func(frame frames.Frame) bool) *FrameFilter {
	return &FrameFilter{
		filterFn: filterFn,
	}
}

// isControlFrame checks if a frame is a control frame that should always pass through.
func (f *FrameFilter) isControlFrame(frame frames.Frame) bool {
	switch frame.(type) {
	case *frames.StartFrame, *frames.EndFrame, *frames.CancelFrame, *frames.ErrorFrame, *frames.StopTaskFrame, *frames.MetricsFrame, *frames.SyncFrame:
		return true
	default:
		return false
	}
}

// ProcessFrame applies the filter function and pushes the frame if it passes.
func (f *FrameFilter) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Always pass control frames through.
	if f.isControlFrame(frame) {
		f.PushFrame(frame, direction)
		return
	}

	// Apply the filter to data frames.
	if f.filterFn(frame) {
		f.PushFrame(frame, direction)
	}
}
