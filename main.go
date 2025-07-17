package main

import (
	"log"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
	"github.com/wuyong/pipeline-go/pkg/processors/filters"
)

func main() {
	log.Println("Starting null filter example")

	// 1. Create a null filter and a logger
	filter := filters.NewNullFilter()
	logger := processors.NewLoggerProcessor("OutputLogger")

	// 2. Create a pipeline
	pl := pipeline.NewPipeline(
		[]processors.FrameProcessor{filter, logger},
		nil, nil,
	)

	// 3. Create a pipeline task
	params := pipeline.PipelineParams{}
	task := pipeline.NewPipelineTask(pl, params)

	// 4. Run the task in a separate goroutine
	go task.Run()

	// 5. Queue some frames
	log.Println("Queueing frames")
	task.QueueFrame(&frames.TextFrame{Text: "this should be dropped by the null filter"})
	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Null filter example finished")
}
