package pipeline

import (
	"sync"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// ParallelPipeline runs multiple pipelines in parallel.
type ParallelPipeline struct {
	processors.BaseProcessor
	pipelines []processors.FrameProcessor
	upQueue   chan frames.Frame
	downQueue chan frames.Frame
	wg        sync.WaitGroup
	once      sync.Once
}

// NewParallelPipeline creates a new ParallelPipeline.
// Each argument is a list of processors for a single parallel pipeline.
func NewParallelPipeline(pipelines ...[]processors.FrameProcessor) *ParallelPipeline {
	if len(pipelines) == 0 {
		panic("ParallelPipeline needs at least one pipeline")
	}

	pp := &ParallelPipeline{
		upQueue:   make(chan frames.Frame, 128),
		downQueue: make(chan frames.Frame, 128),
	}

	for _, procs := range pipelines {
		// Create a standard pipeline for the list of processors.
		// Wire up the pipeline's source and sink to the parallel pipeline's queues.
		pipeline := NewPipeline(procs,
			func(frame frames.Frame, direction processors.FrameDirection) {
				pp.upQueue <- frame
			},
			func(frame frames.Frame, direction processors.FrameDirection) {
				pp.downQueue <- frame
			},
		)
		pp.pipelines = append(pp.pipelines, pipeline)
	}

	return pp
}

// ProcessFrame fans out frames to all parallel pipelines.
func (pp *ParallelPipeline) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Start the queue processing goroutines only once.
	pp.once.Do(func() {
		pp.startQueueProcessors()
	})

	// On cancellation or end, close the queues to terminate the goroutines.
	switch frame.(type) {
	case frames.CancelFrame, frames.EndFrame:
		// Closing the channel will terminate the range loop in the fan-in goroutines
		// which will then cause the WaitGroup to be done.
		defer close(pp.downQueue)
		defer close(pp.upQueue)
	}

	// Fan out the frame to all pipelines.
	// Each pipeline will run in its own goroutine to process the frame.
	var wg sync.WaitGroup
	for _, p := range pp.pipelines {
		wg.Add(1)
		go func(pl processors.FrameProcessor) {
			defer wg.Done()
			pl.ProcessFrame(frame, direction)
		}(p)
	}
	wg.Wait()
}

// startQueueProcessors starts the goroutines that fan-in results from the parallel pipelines.
func (pp *ParallelPipeline) startQueueProcessors() {
	pp.wg.Add(2)
	go pp.processUpQueue()
	go pp.processDownQueue()
}

// processUpQueue collects frames from the upstream and pushes them out.
func (pp *ParallelPipeline) processUpQueue() {
	defer pp.wg.Done()
	for frame := range pp.upQueue {
		pp.PushFrame(frame, processors.FrameDirectionUpstream)
	}
}

// processDownQueue collects frames from the downstream and pushes them out.
func (pp *ParallelPipeline) processDownQueue() {
	defer pp.wg.Done()
	for frame := range pp.downQueue {
		pp.PushFrame(frame, processors.FrameDirectionDownstream)
	}
}

// Cleanup waits for all goroutines to finish.
func (pp *ParallelPipeline) Cleanup() {
	pp.wg.Wait()
	for _, p := range pp.pipelines {
		p.Cleanup()
	}
}
