package main

import (
	"log"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

// SlowProcessor is a processor with an artificial delay.
type SlowProcessor struct {
	processors.BaseProcessor
}

func (p *SlowProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	log.Printf("SlowProcessor: received %T, starting to 'work' for 1 second...", frame)
	time.Sleep(1 * time.Second)
	log.Printf("SlowProcessor: finished 'work' for %T, pushing downstream.", frame)
	p.PushFrame(frame, direction)
}

func main() {
	log.Println("Starting concurrent processor example")

	// 1. Create a slow processor and wrap it in a concurrent one.
	slowProc := &SlowProcessor{}
	concurrentProc := processors.NewConcurrentProcessor(slowProc)
	logger := processors.NewLoggerProcessor("OutputLogger")

	// 2. Create a pipeline with only the concurrent processor
	pl := pipeline.NewPipeline(
		[]processors.FrameProcessor{concurrentProc},
		nil, nil,
	)
	// The logger is linked to the output of the pipeline.
	pl.Link(logger)

	// 3. Create a pipeline task
	params := pipeline.PipelineParams{}
	task := pipeline.NewPipelineTask(pl, params)

	// 4. Run the task in a separate goroutine
	go task.Run()

	// 5. Queue frames. This should return quickly, without waiting for the slow processor.
	log.Println("Queueing frame 1 (should not block)")
	task.QueueFrame(&frames.TextFrame{Text: "first"})
	log.Println("Queueing frame 2 (should not block)")
	task.QueueFrame(&frames.TextFrame{Text: "second"})

	// 6. End the pipeline. This will now block until the concurrent processor is finished.
	log.Println("Queueing EndFrame (will block until work is done)")
	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Concurrent processor example finished")
}
