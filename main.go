package main

import (
	"log"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

// CustomData frame
type CustomData struct {
	Data string
}

func main() {
	log.Println("Starting sync parallel pipeline example")

	// 1. Create two pipelines to be synchronized
	p1 := pipeline.NewPipeline(
		[]processors.FrameProcessor{processors.NewLoggerProcessor("P1-Sync")},
		nil, nil,
	)
	p2 := pipeline.NewPipeline(
		[]processors.FrameProcessor{processors.NewLoggerProcessor("P2-Sync")},
		nil, nil,
	)

	// 2. Create a sync parallel pipeline
	// This pipeline will also have a logger attached to its output.
	spp := pipeline.NewSyncParallelPipeline(p1, p2)
	spp.Link(processors.NewLoggerProcessor("SPP-Output"))

	// 3. Create a pipeline task
	params := pipeline.PipelineParams{}
	task := pipeline.NewPipelineTask(spp, params)

	// 4. Run the task in a separate goroutine
	go task.Run()

	// 5. Queue some frames
	log.Println("Queueing frames")
	task.QueueFrame(&CustomData{Data: "Hello from Sync"})
	task.QueueFrame(frames.EndFrame{})

	// The task will finish automatically when it processes the EndFrame.
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Sync parallel pipeline example finished")
}