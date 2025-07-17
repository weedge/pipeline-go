package main

import (
	"log"
	"reflect"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/notifiers"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
	"github.com/wuyong/pipeline-go/pkg/processors/aggregators"
)

func main() {
	log.Println("Starting notifier-based hold aggregator example")

	// 1. Create a notifier, aggregator, and logger
	notifier := notifiers.NewChannelNotifier()
	aggregator := aggregators.NewHoldLastFrameAggregator(reflect.TypeOf(&frames.TextFrame{}), notifier)
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

	// 5. Queue a frame to be held
	log.Println("Queueing frame to hold")
	task.QueueFrame(&frames.TextFrame{Text: "this frame is held"})

	// 6. Start a goroutine to notify after a delay
	go func() {
		log.Println("Notifier will send notification in 1 second...")
		time.Sleep(1 * time.Second)
		log.Println("Notifying!")
		notifier.Notify()
	}()

	// 7. End the pipeline after giving the notification time to be processed
	time.Sleep(1500 * time.Millisecond)
	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Notifier example finished")
}
