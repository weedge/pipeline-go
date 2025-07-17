package main

import (
	"log"
	"strings"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
	"github.com/wuyong/pipeline-go/pkg/processors/filters"
)

func main() {
	log.Println("Starting filter example")

	// 1. Create a filter and a logger
	filter := filters.NewFrameFilter(func(frame frames.Frame) bool {
		if tf, ok := frame.(*frames.TextFrame); ok {
			return strings.Contains(tf.Text, "keep")
		}
		return false
	})
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
	task.QueueFrame(&frames.TextFrame{Text: "this one should be dropped"})
	task.QueueFrame(&frames.TextFrame{Text: "this one we should keep"})
	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Filter example finished")
}
