package main

import (
	"log/slog"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
	"github.com/weedge/pipeline-go/pkg/processors/filters"
)

func main() {
	slog.Info("Starting pipeline example...")

	// 1. Create the processo
	// This filter will only allow frames that are NOT ImageRawFrames to pass
	imageFilter := filters.NewFrameFilter(func(frame frames.Frame) bool {
		if _, ok := frame.(*frames.ImageRawFrame); ok {
			return false // Drop the frame
		}
		return true // Keep the frame
	})

	// This logger will print any frame that it receives
	logger := processors.NewFrameTraceLogger("output", 0)

	// 2. Create a new pipeline
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			imageFilter,
			logger,
		},
		nil, nil,
	)

	// 3. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	go task.Run()

	// 4. Send frames to the pipeline
	slog.Info("Queueing TextFrame...")
	task.QueueFrame(frames.NewTextFrame("Hello, world!"))

	slog.Info("Queueing ImageRawFrame (will be filtered out)...")
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{}, "PNG", "RGB"))

	slog.Info("Queueing AudioRawFrame...")
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))

	// 5. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewStopTaskFrame())
	time.Sleep(1 * time.Second)

	slog.Info("Pipeline finished.")
}
