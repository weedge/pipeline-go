package main

import (
	"log"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
	"github.com/weedge/pipeline-go/pkg/processors/filters"
)

func main() {
	log.Println("Starting pipeline example...")

	// 1. Create the processors
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
		[]processors.FrameProcessor{
			imageFilter,
			logger,
		},
		nil, nil,
	)

	// 3. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	go task.Run()

	// 4. Send frames to the pipeline
	log.Println("Queueing TextFrame...")
	task.QueueFrame(frames.NewTextFrame("Hello, world!"))

	log.Println("Queueing ImageRawFrame (will be filtered out)...")
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{}, "PNG", "RGB"))

	log.Println("Queueing AudioRawFrame...")
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))

	// 5. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewStopTaskFrame())

	log.Println("Pipeline finished.")
}