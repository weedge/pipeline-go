package pipeline

import (
	"fmt"
	"strings"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// PipelineSource is the entry point for frames into the pipeline.
type PipelineSource struct {
	processors.BaseProcessor
	upstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)
}

func NewPipelineSource(upstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)) *PipelineSource {
	return &PipelineSource{
		upstreamPushFrame: upstreamPushFrame,
	}
}

func (s *PipelineSource) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	//log.Printf("PipelineSource %T, %d", frame, direction)
	switch direction {
	case processors.FrameDirectionUpstream:
		s.upstreamPushFrame(frame, direction)
	case processors.FrameDirectionDownstream:
		s.PushFrame(frame, direction)
	}
}

// PipelineSink is the exit point for frames from the pipeline.
type PipelineSink struct {
	processors.BaseProcessor
	downstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)
}

func NewPipelineSink(downstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)) *PipelineSink {
	return &PipelineSink{
		downstreamPushFrame: downstreamPushFrame,
	}
}

func (s *PipelineSink) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	//log.Printf("PipelineSink %T, %d", frame, direction)
	switch direction {
	case processors.FrameDirectionUpstream:
		s.PushFrame(frame, direction)
	case processors.FrameDirectionDownstream:
		s.downstreamPushFrame(frame, direction)
	}
}

// Pipeline is a sequence of FrameProcessors.
type Pipeline struct {
	processors.BaseProcessor
	processors []processors.FrameProcessor
	source     *PipelineSource
	sink       *PipelineSink
}

func NewPipeline(procs []processors.FrameProcessor, up, down func(frames.Frame, processors.FrameDirection)) *Pipeline {
	p := &Pipeline{}

	upPush := p.PushFrame
	if up != nil {
		upPush = up
	}

	downPush := p.PushFrame
	if down != nil {
		downPush = down
	}

	p.source = NewPipelineSource(upPush)
	p.sink = NewPipelineSink(downPush)
	p.processors = append([]processors.FrameProcessor{p.source}, procs...)
	p.processors = append(p.processors, p.sink)
	p.linkProcessors()
	return p
}

func (p *Pipeline) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	if direction == processors.FrameDirectionDownstream {
		p.source.ProcessFrame(frame, processors.FrameDirectionDownstream)
	} else if direction == processors.FrameDirectionUpstream {
		p.sink.ProcessFrame(frame, processors.FrameDirectionUpstream)
	}
}

func (p *Pipeline) Cleanup() {
	for _, proc := range p.processors {
		proc.Cleanup()
	}
}

func (p *Pipeline) linkProcessors() {
	if len(p.processors) == 0 {
		return
	}
	prev := p.processors[0]
	for _, curr := range p.processors[1:] {
		prev.Link(curr)
		prev = curr
	}
}

func (p *Pipeline) String() string {
	var names []string
	for _, proc := range p.processors {
		if proc == p.source {
			names = append(names, "Source")
		} else if proc == p.sink {
			names = append(names, "Sink")
		} else {
			names = append(names, fmt.Sprintf("%T", proc))
		}
	}
	return "Pipeline: " + strings.Join(names, " -> ")
}
