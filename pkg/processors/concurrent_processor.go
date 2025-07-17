package processors

import (
	"sync"

	"github.com/wuyong/pipeline-go/pkg/frames"
)

// ConcurrentProcessor wraps a FrameProcessor to make it run in its own goroutine.
type ConcurrentProcessor struct {
	BaseProcessor
	wrappedProcessor FrameProcessor
	inQueue          chan frames.Frame
	direction        FrameDirection
	once             sync.Once
	wg               sync.WaitGroup
}

// NewConcurrentProcessor creates a new ConcurrentProcessor.
func NewConcurrentProcessor(processor FrameProcessor) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		wrappedProcessor: processor,
		inQueue:          make(chan frames.Frame, 128),
	}
}

// startWorker starts the internal goroutine that processes frames from the queue.
func (p *ConcurrentProcessor) startWorker() {
	p.wrappedProcessor.Link(&p.BaseProcessor)
	go func() {
		for frame := range p.inQueue {
			p.wrappedProcessor.ProcessFrame(frame, p.direction)
			p.wg.Done() // Signal that one frame is done processing.
		}
	}()
}

// ProcessFrame puts the frame onto the internal queue and manages the WaitGroup.
func (p *ConcurrentProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	p.once.Do(func() {
		p.direction = direction
		p.startWorker()
	})

	if _, ok := frame.(frames.EndFrame); ok {
		p.wg.Wait()      // Wait for all queued frames to finish.
		close(p.inQueue) // Close the queue to terminate the worker.
		// After the worker is done, we push the EndFrame ourselves.
		p.PushFrame(frame, direction)
	} else {
		p.wg.Add(1) // Increment the counter for each new frame.
		p.inQueue <- frame
	}
}
