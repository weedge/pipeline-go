package main

import (
	"log"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func main() {
	log.Println("Starting async pipeline example...")

	// 1. Create an async processor
	asyncProc := processors.NewAsyncFrameProcessor("async_processor")
	asyncProc.SetVerbose(true)

	// 2. Link it to a logger processor
	logger := processors.NewFrameTraceLogger("async", 0)
	logger.SetVerbose(true)

	// 3. Create a simple pipeline with the async processor
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			asyncProc,
			logger,
		},
		nil, nil,
	)
	log.Println(myPipeline)

	// 4. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	go task.Run()

	// 5. Send frames to the pipeline
	log.Println("Queueing frames...")
	task.QueueFrame(frames.NewTextFrame("Hello, async world!"))
	task.QueueFrame(frames.NewTextFrame("Processing asynchronously..."))

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	log.Println("Send a stop frame to terminate the pipeline.")
	// 6. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewEndFrame())

	log.Println("Terminating the pipeline.")
	// Give some time for termination
	time.Sleep(100 * time.Millisecond)

	log.Println("Async pipeline finished.")
}
