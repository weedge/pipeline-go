package main

import (
	"log"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
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

// CanGenerateMetrics returns whether the processor can generate metrics
func (p *CustomProcessor) CanGenerateMetrics() bool {
	return true
}

// ProcessFrame implements the IFrameProcessor interface
func (p *CustomProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	log.Printf("Processing frame %s with direction %s\n", frame.Name(), direction.String())
	// Start metrics collection
	p.StartProcessingMetrics()
	defer p.StopProcessingMetrics()

	// Call the parent implementation
	p.FrameProcessor.ProcessFrame(frame, direction)
	p.PushFrame(frame, direction)
}

func main() {
	log.Println("Starting metrics example...")

	// 1. Create a custom processor that can generate metrics
	customProc := NewCustomProcessor("custom-processor")
	//customProc.SetVerbose(true)

	// 2. Link it to a logger processor
	logger := processors.NewFrameTraceLogger("metrics", 0) // Add small delay to see the processing

	// 3. Create a simple pipeline with the custom processor
	myPipeline := pipeline.NewPipelineWithVerbose(
		[]processors.IFrameProcessor{
			customProc,
			logger,
		},
		nil, nil, true,
	)

	// 4. Create and run a pipeline task with metrics enabled to push start frame
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{
		EnableMetrics: true,
	})
	go task.Run()

	// 6. Send frames to the pipeline
	log.Println("Queueing frames...")
	task.QueueFrame(frames.NewTextFrame("Hello, metrics world!"))
	task.QueueFrame(frames.NewTextFrame("Processing with metrics..."))

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	log.Println("Send a stop frame to terminate the pipeline.")
	// 8. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewEndFrame())

	log.Println("Terminating the pipeline.")
	// Give some time for termination
	time.Sleep(100 * time.Millisecond)

	log.Println("Metrics example finished.")
}
