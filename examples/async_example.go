package main

import (
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func main() {
	logger.InitLoggerWithConfig(logger.NewDefaultLoggerConfig())
	logger.Info("Starting async pipeline example...")

	// 1. Create an async processor
	asyncProc := processors.NewAsyncFrameProcessor("async_processor")
	asyncProc.WithPorcessFrameAllowPush(true).SetVerbose(true)

	// 2. Link it to a logger processor
	trace_logger := processors.NewFrameTraceLogger("async", 0)
	trace_logger.SetVerbose(true)

	// 3. Create a simple pipeline with the async processor
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			asyncProc,
			trace_logger,
		},
		nil, nil,
	)
	logger.Info(myPipeline.String())

	// 4. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	go task.Run()

	// 5. Send frames to the pipeline
	logger.Info("Queueing frames...")
	task.QueueFrame(frames.NewTextFrame("Hello, async world!"))
	task.QueueFrame(frames.NewTextFrame("Processing asynchronously..."))

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	logger.Info("Send a stop frame to terminate the pipeline.")
	// 6. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewEndFrame())

	logger.Info("Terminating the pipeline.")
	// Give some time for termination
	time.Sleep(100 * time.Millisecond)

	logger.Info("Async pipeline finished.")
}
