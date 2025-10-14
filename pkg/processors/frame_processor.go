package processors

import (
	"fmt"
	"runtime"
	"slices"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
)

// FrameDirection indicates the direction of a frame in the pipeline.
type FrameDirection int

const (
	// FrameDirectionUpstream indicates a frame is moving from sink to source.
	FrameDirectionUpstream FrameDirection = iota
	// FrameDirectionDownstream indicates a frame is moving from source to sink.
	FrameDirectionDownstream
)

// String returns the string representation of FrameDirection
func (d FrameDirection) String() string {
	switch d {
	case FrameDirectionUpstream:
		return "Upstream(<-)"
	case FrameDirectionDownstream:
		return "Downstream(->)"
	default:
		return "Unknown"
	}
}

// IFrameProcessor is the interface for all components that process frames.
type IFrameProcessor interface {
	Name() string
	Link(p IFrameProcessor)
	SetPrev(prev IFrameProcessor)
	ProcessFrame(frame frames.Frame, direction FrameDirection)
	Cleanup()
	SetVerbose(verbose bool)
}

// FrameProcessor provides a base implementation for IFrameProcessor.
// with metrics, error handling, and other advanced features.
type FrameProcessor struct {
	id                    int64
	name                  string
	parentPipeline        IFrameProcessor
	next                  IFrameProcessor
	prev                  IFrameProcessor
	allowInterruptions    bool
	enableMetrics         bool
	enableUsageMetrics    bool
	reportOnlyInitialTTFB bool
	metrics               *MetricsProcessor
	skipFrames            []frames.Frame
	verbose               bool
}

// NewFrameProcessor creates a new FrameProcessor.
func NewFrameProcessor(name string) *FrameProcessor {
	return &FrameProcessor{
		name:       name,
		metrics:    NewMetricsProcessor(name),
		skipFrames: make([]frames.Frame, 0),
	}
}

// ID returns the processor's ID.
func (p *FrameProcessor) ID() int64 {
	return p.id
}

// Name returns the processor's name.
func (p *FrameProcessor) Name() string {
	return p.name
}

// InterruptionsAllowed returns whether interruptions are allowed.
func (p *FrameProcessor) InterruptionsAllowed() bool {
	return p.allowInterruptions
}

// MetricsEnabled returns whether metrics are enabled.
func (p *FrameProcessor) MetricsEnabled() bool {
	return p.enableMetrics
}

// UsageMetricsEnabled returns whether usage metrics are enabled.
func (p *FrameProcessor) UsageMetricsEnabled() bool {
	return p.enableUsageMetrics
}

// ReportOnlyInitialTTFB returns whether only initial TTFB should be reported.
func (p *FrameProcessor) ReportOnlyInitialTTFB() bool {
	return p.reportOnlyInitialTTFB
}

// Verbose returns whether verbose mode is enabled.
func (p *FrameProcessor) Verbose() bool {
	return p.verbose
}

func (p *FrameProcessor) SetVerbose(verbose bool) {
	p.verbose = verbose
}

// AddSkipFrame adds a frame to the skip list.
func (p *FrameProcessor) AddSkipFrame(frame frames.Frame) {
	p.skipFrames = append(p.skipFrames, frame)
}

// StartTTFBMetrics starts TTFB metrics collection.
func (p *FrameProcessor) StartTTFBMetrics() {
	if p.MetricsEnabled() {
		p.metrics.StartTTFBMetrics(p.ReportOnlyInitialTTFB())
	}
}

// StopTTFBMetrics stops TTFB metrics collection and pushes a MetricsFrame if applicable.
func (p *FrameProcessor) StopTTFBMetrics() {
	if p.MetricsEnabled() {
		frame := p.metrics.StopTTFBMetrics()
		if frame != nil {
			p.PushFrame(frame, FrameDirectionDownstream)
		}
	}
}

// StartProcessingMetrics starts processing time metrics collection.
func (p *FrameProcessor) StartProcessingMetrics() {
	if p.MetricsEnabled() {
		p.metrics.StartProcessingMetrics()
	}
}

// StopProcessingMetrics stops processing time metrics collection and pushes a MetricsFrame if applicable.
func (p *FrameProcessor) StopProcessingMetrics() {
	if p.MetricsEnabled() {
		frame := p.metrics.StopProcessingMetrics()
		if frame != nil {
			p.PushFrame(frame, FrameDirectionDownstream)
		}
	}
}

// StopAllMetrics stops all metrics collection.
func (p *FrameProcessor) StopAllMetrics() {
	p.StopTTFBMetrics()
	p.StopProcessingMetrics()
}

// SetParentPipeline sets the parent pipeline.
func (p *FrameProcessor) SetParentPipeline(pipeline IFrameProcessor) {
	p.parentPipeline = pipeline
}

// GetParentPipeline returns the parent pipeline.
func (p *FrameProcessor) GetParentPipeline() IFrameProcessor {
	return p.parentPipeline
}

// ProcessFrame implements the IFrameProcessor interface.
// Handle StartFrame to init and StartInterruptionFrame to stop all metrics
func (p *FrameProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	// Check if frame should be skipped
	if slices.Contains(p.skipFrames, frame) {
		return
	}

	// Handle StartFrame
	if startFrame, ok := frame.(*frames.StartFrame); ok {
		p.allowInterruptions = startFrame.AllowInterruptions
		p.enableMetrics = startFrame.EnableMetrics
		p.enableUsageMetrics = startFrame.EnableUsageMetrics
		p.reportOnlyInitialTTFB = startFrame.ReportOnlyInitialTTFB
	} else if _, ok := frame.(*frames.StartInterruptionFrame); ok {
		// Handle StartInterruptionFrame
		p.StopAllMetrics()
	}

}

// PushError pushes an error frame upstream.
func (p *FrameProcessor) PushError(errorFrame *frames.ErrorFrame) {
	p.PushFrame(errorFrame, FrameDirectionUpstream)
}

// PushDownstreamFrame pushes a frame in the Downstream direction.
func (p *FrameProcessor) PushDownstreamFrame(frame frames.Frame) {
	p.PushFrame(frame, FrameDirectionDownstream)
}

// PushUpstreamFrame pushes a frame in the Upstream direction.
func (p *FrameProcessor) PushUpstreamFrame(frame frames.Frame) {
	p.PushFrame(frame, FrameDirectionUpstream)
}

// PushFrame pushes a frame in the specified direction.
func (p *FrameProcessor) PushFrame(frame frames.Frame, direction FrameDirection) {
	defer func() {
		if r := recover(); r != nil {
			// 获取调用栈信息
			buf := make([]byte, 10240)
			n := runtime.Stack(buf, false)
			stackTrace := string(buf[:n])
			println(stackTrace) // print in stderr for panic stack

			logger.Info(fmt.Sprintf("Uncaught panic in %s(%T): %v\nStack trace:\n%s", p.name, p, r, stackTrace))
		}
	}()

	if direction == FrameDirectionDownstream && p.next != nil {
		if p.verbose {
			logger.Info(fmt.Sprintf("Downstream %d Pushing %s  %s(%T) -> %s (Calling ProcessFrame on next: %T)", direction, frame.String(), p.name, p, p.next.Name(), p.next))
		}
		p.next.ProcessFrame(frame, direction)
	} else if direction == FrameDirectionUpstream && p.prev != nil {
		if p.verbose {
			logger.Info(fmt.Sprintf("Upstream %d Pushing %s  %s(%T) -> %s (Calling ProcessFrame on prev: %T)", direction, frame.String(), p.name, p, p.prev.Name(), p.prev))
		}

		// Check if prev is a mockProcessor with a custom prevProcessor field
		// We need to use reflection to check for this
		p.prev.ProcessFrame(frame, direction)
	} else {
		if p.verbose {
			logger.Info(fmt.Sprintf("Frame not pushed: direction=%d, next=%v, prev=%v", direction, p.next, p.prev))
		}
	}
}

// Link implements the IFrameProcessor interface.
func (p *FrameProcessor) Link(next IFrameProcessor) {
	p.next = next
	if p.verbose {
		logger.Info(fmt.Sprintf("%s(%T) -> %s(%T)", p.Name(), p, next.Name(), next))
	}
	//next.SetPrev(p) #ISSUE: don't use this,  pre FrameProcessor,not caller
}

// SetPrev sets the previous processor in the pipeline.
func (p *FrameProcessor) SetPrev(prev IFrameProcessor) {
	if p.verbose {
		logger.Info(fmt.Sprintf("%s(%T) <- %s(%T)", prev.Name(), prev, p.Name(), p))
	}
	p.prev = prev
}

func (p *FrameProcessor) Cleanup() {
	// Base implementation can be empty.
}

// String returns the processor's name.
func (p *FrameProcessor) String() string {
	return p.name
}
