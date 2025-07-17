package aggregators

import (
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// GatedAggregator accumulates frames and releases them when a gate opens.
type GatedAggregator struct {
	processors.BaseProcessor
	gateOpenFn  func(frames.Frame) bool
	gateCloseFn func(frames.Frame) bool
	isGateOpen  bool
	accumulator []frames.Frame
	direction   processors.FrameDirection
}

// NewGatedAggregator creates a new GatedAggregator.
func NewGatedAggregator(
	gateOpenFn, gateCloseFn func(frames.Frame) bool,
	startOpen bool,
	direction processors.FrameDirection,
) *GatedAggregator {
	return &GatedAggregator{
		gateOpenFn:  gateOpenFn,
		gateCloseFn: gateCloseFn,
		isGateOpen:  startOpen,
		accumulator: make([]frames.Frame, 0),
		direction:   direction,
	}
}

// isControlFrame checks if a frame is a control frame that should always pass through.
func (a *GatedAggregator) isControlFrame(frame frames.Frame) bool {
	switch frame.(type) {
	case frames.StartFrame, frames.EndFrame, frames.CancelFrame, frames.ErrorFrame, frames.StopTaskFrame, frames.MetricsFrame, frames.SyncFrame:
		return true
	default:
		return false
	}
}

func (a *GatedAggregator) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Always pass control frames through.
	if a.isControlFrame(frame) {
		a.PushFrame(frame, direction)
		return
	}

	// Ignore frames not matching the configured direction.
	if direction != a.direction {
		a.PushFrame(frame, direction)
		return
	}

	// First, process the frame based on the current gate state.
	if a.isGateOpen {
		a.PushFrame(frame, direction)
	}

	// Then, update the gate state for the next frame.
	if a.isGateOpen {
		if a.gateCloseFn(frame) {
			a.isGateOpen = false
		}
	} else {
		if a.gateOpenFn(frame) {
			a.isGateOpen = true
		}
	}
}
