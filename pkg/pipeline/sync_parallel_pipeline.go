package pipeline

import (
	"sync"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// SyncParallelPipeline runs multiple pipelines in parallel and synchronizes their output.
type SyncParallelPipeline struct {
	processors.FrameProcessor
	pipelines []processors.IFrameProcessor
}

// NewSyncParallelPipeline creates a new SyncParallelPipeline.
func NewSyncParallelPipeline(pipelines ...processors.IFrameProcessor) *SyncParallelPipeline {
	if len(pipelines) == 0 {
		panic("SyncParallelPipeline needs at least one pipeline")
	}

	return &SyncParallelPipeline{
		pipelines: pipelines,
	}
}

// ProcessFrame fans out frames and waits for all pipelines to return a result.
func (spp *SyncParallelPipeline) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	var wg sync.WaitGroup
	results := make(chan frames.Frame, len(spp.pipelines))

	// Fan out: Start a goroutine for each pipeline.
	for _, p := range spp.pipelines {
		wg.Add(1)
		go func(pl processors.IFrameProcessor) {
			defer wg.Done()

			// Create a temporary processor to intercept the result from the pipeline.
			var resultFrame frames.Frame
			var resultMux sync.Mutex
			var resultCond = sync.NewCond(&resultMux)
			isDone := false

			tempProcessor := &resultInterceptor{
				onFrame: func(f frames.Frame, d processors.FrameDirection) {
					resultMux.Lock()
					defer resultMux.Unlock()
					if !isDone {
						resultFrame = f
						isDone = true
						resultCond.Signal()
					}
				},
			}
			pl.Link(tempProcessor)

			// Send the frame to the pipeline.
			pl.ProcessFrame(frame, direction)

			// Wait for the result.
			resultMux.Lock()
			for !isDone {
				resultCond.Wait()
			}
			resultMux.Unlock()

			// Send the result to the results channel.
			if resultFrame != nil {
				results <- resultFrame
			}
		}(p)
	}

	// Wait for all pipelines to finish processing the frame.
	wg.Wait()
	close(results)

	// Fan in: Collect all results and push them downstream.
	// Here we are pushing them one by one. A different strategy could be to
	// aggregate them into a single slice frame.
	for res := range results {
		spp.PushFrame(res, direction)
	}
}

// resultInterceptor is a helper processor to capture the output of a pipeline.
type resultInterceptor struct {
	processors.FrameProcessor
	onFrame func(frame frames.Frame, direction processors.FrameDirection)
}

func (ri *resultInterceptor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	ri.onFrame(frame, direction)
}
