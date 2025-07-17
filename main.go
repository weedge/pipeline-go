package main

import (
	"log"
	"reflect"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
	"github.com/wuyong/pipeline-go/pkg/processors/aggregators"
)

func main() {
	log.Println("Starting hold aggregator example")

	// 1. Create a hold aggregator and a logger
	aggregator := aggregators.NewHoldLastFrameAggregator(reflect.TypeOf(&frames.TextFrame{}))
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

	// 5. Queue some frames
	log.Println("Queueing frames")
	task.QueueFrame(&frames.TextFrame{Text: "this is the first frame, it will be replaced"})
	task.QueueFrame(&frames.TextFrame{Text: "this is the last frame, it will be held"})

	// 6. Wait a moment, then release the held frame
	log.Println("Waiting to release...")
	time.Sleep(500 * time.Millisecond)
	log.Println("Releasing frame!")
	aggregator.Release()
	time.Sleep(100 * time.Millisecond) // Give time for the released frame to be processed.

	// 7. End the pipeline
	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Hold aggregator example finished")
}
