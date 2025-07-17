package main

import (
	"log"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

func main() {
	log.Println("Starting idle processor example")

	// 1. Create an idle processor and a logger
	idle := processors.NewIdleProcessor(500 * time.Millisecond)
	logger := processors.NewLoggerProcessor("OutputLogger")

	// 2. Create a pipeline
	pl := pipeline.NewPipeline(
		[]processors.FrameProcessor{idle, logger},
		nil, nil,
	)

	// 3. Create a pipeline task
	params := pipeline.PipelineParams{}
	task := pipeline.NewPipelineTask(pl, params)

	// 4. Run the task in a separate goroutine
	go task.Run()

	// 5. Queue some frames with delays
	log.Println("Queueing frame 1")
	task.QueueFrame(&frames.TextFrame{Text: "first frame"})
	time.Sleep(200 * time.Millisecond)

	log.Println("Queueing frame 2")
	task.QueueFrame(&frames.TextFrame{Text: "second frame"})

	log.Println("Waiting for idle timeout...")
	time.Sleep(700 * time.Millisecond) // Wait for idle frame

	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Idle processor example finished")
}
