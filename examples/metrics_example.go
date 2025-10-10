package main

import (
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// CustomProcessor is a custom processor that can generate metrics
type CustomProcessor struct {
	*processors.FrameProcessor
}

// NewCustomProcessor creates a new CustomProcessor
func NewCustomProcessor(name string) *CustomProcessor {
	return &CustomProcessor{
		FrameProcessor: processors.NewFrameProcessor(name),
	}
}

// ProcessFrame implements the IFrameProcessor interface
func (p *CustomProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// First Call the parent implementation to init with start frame
	p.FrameProcessor.ProcessFrame(frame, direction)

	// Start metrics collection
	p.StartProcessingMetrics()
	defer p.StopProcessingMetrics()

	p.PushFrame(frame, direction)

	// mock process time
	switch frame.(type) {
	case *frames.StopTaskFrame, *frames.EndFrame:
		time.Sleep(1 * time.Millisecond)
		return
	default:
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	logger.Info("Starting metrics example...")

	// 1. Create a custom processor that can generate metrics
	customProc := NewCustomProcessor("custom-processor")
	//customProc.SetVerbose(true)

	// 2. Link it to a logger processor
	trace_logger := processors.NewDefaultFrameLoggerProcessorWithName("metrics")

	// 3. Create a simple pipeline with the custom processor
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			customProc,
			trace_logger,
		},
		nil, nil,
	)
	logger.Info(myPipeline.String())

	// 4. Create and run a pipeline task with metrics enabled to push start frame
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{
		EnableMetrics: true,
	})
	go task.Run()

	// 6. Send frames to the pipeline
	logger.Info("Queueing frames...")
	task.QueueFrame(frames.NewTextFrame("Hello, metrics world!"))
	task.QueueFrame(frames.NewTextFrame("Processing with metrics..."))

	// Give some time for async processing
	time.Sleep(1000 * time.Millisecond)

	logger.Info("Send a stop frame to terminate the pipeline.")
	// 8. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewEndFrame())

	logger.Info("Terminating the pipeline.")
	// Give some time for termination
	time.Sleep(100 * time.Millisecond)

	logger.Info("Metrics example finished.")
}
