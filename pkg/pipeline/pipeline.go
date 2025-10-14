package pipeline

import (
	"fmt"
	"strings"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// PipelineSource is the entry point for frames into the pipeline.
type PipelineSource struct {
	processors.FrameProcessor
	upstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)
}

func NewPipelineSource(upstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)) *PipelineSource {
	return &PipelineSource{ // new a default FrameProcessor
		upstreamPushFrame: upstreamPushFrame,
	}
}

func (s *PipelineSource) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	//logger.Infof("PipelineSource %T, %d", frame, direction)
	switch frame.(type) {
	case *frames.StartFrame:
		logger.Infof("PipelineSource get %T direction %s Pipeline Start!", frame, direction.String())
	case *frames.EndFrame:
		logger.Infof("PipelineSource get %T direction %s Pipeline End!", frame, direction.String())
	case *frames.CancelFrame:
		logger.Infof("PipelineSource get %T direction %s Pipeline Cancel!", frame, direction.String())
	}

	switch direction {
	case processors.FrameDirectionUpstream:
		s.upstreamPushFrame(frame, direction)
	case processors.FrameDirectionDownstream:
		s.PushFrame(frame, direction)
	}
}

func (s *PipelineSource) Name() string {
	return "PipelineSource"
}

// PipelineSink is the exit point for frames from the pipeline.
type PipelineSink struct {
	processors.FrameProcessor
	downstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)
}

func NewPipelineSink(downstreamPushFrame func(frame frames.Frame, direction processors.FrameDirection)) *PipelineSink {
	return &PipelineSink{
		downstreamPushFrame: downstreamPushFrame,
	}
}

func (s *PipelineSink) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	//logger.Infof("PipelineSink %T, %d", frame, direction)
	switch frame.(type) {
	case *frames.StartFrame:
		logger.Infof("PipelineSink get %T direction %s Pipeline Start!", frame, direction.String())
	case *frames.EndFrame:
		logger.Infof("PipelineSink get %T direction %s Pipeline End!", frame, direction.String())
	case *frames.CancelFrame:
		logger.Infof("PipelineSink get %T direction %s Pipeline Cancel!", frame, direction.String())
	}

	switch direction {
	case processors.FrameDirectionUpstream:
		s.PushFrame(frame, direction)
	case processors.FrameDirectionDownstream:
		s.downstreamPushFrame(frame, direction)
	}
}

func (s *PipelineSink) Name() string {
	return "PipelineSink"
}

// Pipeline is a sequence of FrameProcessors.
type Pipeline struct {
	processors.FrameProcessor
	processors []processors.IFrameProcessor
	source     *PipelineSource
	sink       *PipelineSink
}

func NewPipeline(procs []processors.IFrameProcessor, up, down func(frames.Frame, processors.FrameDirection)) *Pipeline {
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
	p.processors = append([]processors.IFrameProcessor{p.source}, procs...)
	p.processors = append(p.processors, p.sink)
	p.linkProcessors()
	return p
}

func NewPipelineWithVerbose(procs []processors.IFrameProcessor, up, down func(frames.Frame, processors.FrameDirection), verbose bool) *Pipeline {
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
	p.processors = append([]processors.IFrameProcessor{p.source}, procs...)
	p.processors = append(p.processors, p.sink)
	p.linkProcessorsSetVerbose(verbose)
	return p
}

func (p *Pipeline) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	switch direction {
	case processors.FrameDirectionDownstream:
		p.source.ProcessFrame(frame, processors.FrameDirectionDownstream)
	case processors.FrameDirectionUpstream:
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
		curr.SetPrev(prev)
		prev = curr
	}
}

func (p *Pipeline) linkProcessorsSetVerbose(verbose bool) {
	if len(p.processors) == 0 {
		return
	}
	prev := p.processors[0]
	prev.SetVerbose(verbose)
	for _, curr := range p.processors[1:] {
		curr.SetVerbose(verbose)
		prev.Link(curr)
		curr.SetPrev(prev)
		prev = curr
	}
}

func (p *Pipeline) String() string {
	var names []string
	for _, proc := range p.processors {
		switch proc {
		case p.source:
			names = append(names, "Source")
		case p.sink:
			names = append(names, "Sink")
		default:
			names = append(names, fmt.Sprintf("%T", proc))
		}
	}
	return "Pipeline: " + strings.Join(names, " <-> ")
}
