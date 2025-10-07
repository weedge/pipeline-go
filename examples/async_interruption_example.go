package main

import (
	"log"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func main() {
	log.Println("Starting async interruption example...")

	// 1. Create an async processor
	asyncProc := processors.NewAsyncFrameProcessor("interruption-example")

	// 2. Link it to a logger processor
	logger := processors.NewFrameTraceLogger("async", 10) // Add small delay to see the processing

	// 3. Create a simple pipeline with the async processor
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			asyncProc,
			logger,
		},
		nil, nil,
	)
	log.Println(myPipeline.String())

	// 4. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	go task.Run()

	// 5. Send frames to the pipeline
	log.Println("Queueing frames...")
	task.QueueFrame(frames.NewTextFrame("Hello, async world!"))
	task.QueueFrame(frames.NewTextFrame("Processing asynchronously..."))

	// Wait a bit to let the frames process
	time.Sleep(100 * time.Millisecond)

	// 6. Send an interruption frame
	log.Println("Sending interruption frame...")
	task.QueueFrame(frames.NewStartInterruptionFrame())

	// Wait a bit to let the interruption process
	time.Sleep(50 * time.Millisecond)

	// 7. Send more frames after interruption
	log.Println("Queueing frames after interruption...")
	task.QueueFrame(frames.NewTextFrame("After interruption..."))

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	log.Println("Send a stop frame to terminate the pipeline.")
	// 8. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewEndFrame())

	log.Println("Terminating the pipeline.")
	// Give some time for termination
	time.Sleep(100 * time.Millisecond)

	log.Println("Async interruption example finished.")
}
