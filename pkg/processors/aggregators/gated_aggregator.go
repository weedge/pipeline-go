package aggregators

import (
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// GatedAggregator accumulates frames and releases them when a gate opens.
type GatedAggregator struct {
	processors.FrameProcessor
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
	case *frames.StartFrame, *frames.EndFrame, *frames.CancelFrame, *frames.ErrorFrame, *frames.StopTaskFrame, *frames.MetricsFrame, *frames.SyncFrame:
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

	shouldOpen := !a.isGateOpen && a.gateOpenFn(frame)
	shouldClose := a.isGateOpen && a.gateCloseFn(frame)

	// If the gate is currently open, we should process the frame.
	// This includes the frame that is about to close the gate.
	if a.isGateOpen {
		a.PushFrame(frame, direction)
	} else if shouldOpen {
		// If the gate is closed but this frame opens it, we should also process it.
		a.PushFrame(frame, direction)
	}

	// Now, update the state for the next frame.
	if shouldOpen {
		a.isGateOpen = true
	} else if shouldClose {
		a.isGateOpen = false
	}
}
