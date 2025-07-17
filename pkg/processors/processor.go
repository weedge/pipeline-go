package processors

import (
	"github.com/weedge/pipeline-go/pkg/frames"
)

// FrameDirection indicates the direction of a frame in the pipeline.
type FrameDirection int

const (
	// FrameDirectionUpstream indicates a frame is moving from sink to source.
	FrameDirectionUpstream FrameDirection = iota
	// FrameDirectionDownstream indicates a frame is moving from source to sink.
	FrameDirectionDownstream
)

// FrameProcessor is the interface for all components that process frames.
type FrameProcessor interface {
	ProcessFrame(frame frames.Frame, direction FrameDirection)
	Link(p FrameProcessor)
	Cleanup()
}

// BaseProcessor provides a base implementation for FrameProcessor.
type BaseProcessor struct {
	next FrameProcessor
}

func (p *BaseProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	// Base implementation can be empty, or have some common logic.
}

func (p *BaseProcessor) Link(next FrameProcessor) {
	p.next = next
}

func (p *BaseProcessor) Cleanup() {
	// Base implementation can be empty.
}

func (p *BaseProcessor) PushFrame(frame frames.Frame, direction FrameDirection) {
	if p.next != nil {
		p.next.ProcessFrame(frame, direction)
	}
}
