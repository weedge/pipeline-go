package processors

import "github.com/weedge/pipeline-go/pkg/frames"

// OutputProcessor calls a callback function for each frame.
type OutputProcessor struct {
	FrameProcessor
	callback func(frames.Frame)
}

// NewOutputProcessor creates a new OutputProcessor.
func NewOutputProcessor(callback func(frames.Frame)) *OutputProcessor {
	return &OutputProcessor{
		callback: callback,
	}
}

// ProcessFrame calls the callback function.
func (p *OutputProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	p.callback(frame)
	p.PushFrame(frame, direction)
}
