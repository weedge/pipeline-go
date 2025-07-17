package pipeline

import (
	"sync"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// MergePipeline concurrently merges the output of multiple pipelines into a single stream.
type MergePipeline struct {
	processors.BaseProcessor
	pipelines []processors.FrameProcessor
	outQueue  chan frames.Frame
	wg        sync.WaitGroup
}

// NewMergePipeline creates a new MergePipeline.
func NewMergePipeline(pipelines ...processors.FrameProcessor) *MergePipeline {
	if len(pipelines) == 0 {
		panic("MergePipeline needs at least one pipeline")
	}

	mp := &MergePipeline{
		pipelines: pipelines,
		outQueue:  make(chan frames.Frame, 128),
	}

	// For each incoming pipeline, we create a goroutine that will forward its
	// output to the merge pipeline's output queue.
	mp.wg.Add(len(pipelines))
	for _, p := range mp.pipelines {
		// We need to set the next processor for each pipeline to be a custom
		// function that pushes to the merge queue.
		p.Link(&merger{mp: mp})
	}

	// Create a goroutine to close the output channel once all input pipelines are done.
	go func() {
		mp.wg.Wait()
		close(mp.outQueue)
	}()

	return mp
}

// merger is a helper processor to link input pipelines to the merge queue.
type merger struct {
	processors.BaseProcessor
	mp *MergePipeline
}

func (m *merger) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	m.mp.outQueue <- frame
	if _, ok := frame.(frames.EndFrame); ok {
		m.mp.wg.Done()
	}
}

// ProcessFrame fans out the input frame to all pipelines.
func (mp *MergePipeline) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// We process the frame in each pipeline concurrently.
	var wg sync.WaitGroup
	for _, p := range mp.pipelines {
		wg.Add(1)
		go func(pl processors.FrameProcessor) {
			defer wg.Done()
			pl.ProcessFrame(frame, direction)
		}(p)
	}
	wg.Wait()
}

// GetOutput returns the merged output channel.
func (mp *MergePipeline) GetOutput() <-chan frames.Frame {
	return mp.outQueue
}

func (mp *MergePipeline) Cleanup() {
	// The output channel will be closed by the goroutine when all inputs are done.
	// We just need to call cleanup on the child pipelines.
	for _, p := range mp.pipelines {
		p.Cleanup()
	}
}
