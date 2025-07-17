package main

import (
	"log"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
	"github.com/wuyong/pipeline-go/pkg/processors/aggregators"
)

func main() {
	log.Println("Starting improved aggregator example")

	// 1. Create an aggregator and a logger
	aggregator := aggregators.NewSentenceAggregator()
	logger := processors.NewLoggerProcessor("OutputLogger")

	// 2. Create a pipeline
	pl := pipeline.NewPipeline(
		[]processors.FrameProcessor{aggregator, logger},
		nil, nil,
	)

	// 3. Create a pipeline task
	params := pipeline.PipelineParams{}
	task := pipeline.NewPipelineTask(pl, params)

	// 4. Run the task in a separate goroutine
	go task.Run()

	// 5. Queue some frames to be aggregated
	log.Println("Queueing frames")
	task.QueueFrame(&frames.TextFrame{Text: "This is a test from Dr. Smith"})
	task.QueueFrame(&frames.TextFrame{Text: " in the U.S.A."}) // Should not split here
	task.QueueFrame(&frames.TextFrame{Text: " What do you think?"})
	task.QueueFrame(&frames.TextFrame{Text: "I think it's great!"})
	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Improved aggregator example finished")
}
