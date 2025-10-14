package processors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
)

// AsyncFrameProcessor is a processor that handles frames asynchronously using a queue.
type AsyncFrameProcessor struct {
	*FrameProcessor
	ctx                   context.Context
	cancel                context.CancelFunc
	pushQueueSize         int
	pushQueue             chan pushItem
	pushFrameTask         *sync.WaitGroup
	interruptionMu        sync.Mutex
	porcessFrameAllowPush bool
	passText              bool
	passRawAudio          bool
}

// pushItem represents an item in the push queue.
type pushItem struct {
	frame     frames.Frame
	direction FrameDirection
}

// NewAsyncFrameProcessor creates a new AsyncFrameProcessor.
func NewAsyncFrameProcessor(name string) *AsyncFrameProcessor {
	pushQueueSize := 128
	return NewAsyncFrameProcessorWithPushQueueSize(name, pushQueueSize)
}

func NewAsyncFrameProcessorWithPushQueueSize(name string, pushQueueSize int) *AsyncFrameProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	p := &AsyncFrameProcessor{
		FrameProcessor:        NewFrameProcessor(name),
		ctx:                   ctx,
		cancel:                cancel,
		pushQueueSize:         pushQueueSize,
		pushQueue:             make(chan pushItem, pushQueueSize),
		pushFrameTask:         &sync.WaitGroup{},
		porcessFrameAllowPush: false,
		passText:              false,
		passRawAudio:          false,
	}
	p.createPushTask()
	return p
}

func (p *AsyncFrameProcessor) PassText() bool {
	return p.passText
}

func (p *AsyncFrameProcessor) WithPassText(passText bool) *AsyncFrameProcessor {
	p.passText = passText
	return p
}

func (p *AsyncFrameProcessor) PassRawAudio() bool {
	return p.passRawAudio
}

func (p *AsyncFrameProcessor) WithPassRawAudio(passRawAudio bool) *AsyncFrameProcessor {
	p.passRawAudio = passRawAudio
	return p
}

// ProcessFrameAllowPush returns whether verbose mode is enabled.
func (p *AsyncFrameProcessor) ProcessFrameAllowPush() bool {
	return p.porcessFrameAllowPush
}

func (p *AsyncFrameProcessor) WithPorcessFrameAllowPush(porcessFrameAllowPush bool) *AsyncFrameProcessor {
	p.porcessFrameAllowPush = porcessFrameAllowPush
	return p
}

// ProcessFrame implements the IFrameProcessor interface.
func (p *AsyncFrameProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	// Call base implementation if needed
	p.FrameProcessor.ProcessFrame(frame, direction)

	// Handle interruption frames
	switch frame.(type) {
	case *frames.EndFrame:
		p.Cleanup()
	case *frames.CancelFrame:
		p.Cleanup()
	case *frames.StartInterruptionFrame, frames.StartInterruptionFrame:
		p.HandleInterruptions(frame)
	}

	if p.porcessFrameAllowPush {
		p.QueueFrame(frame, direction)
	}
}

// Cleanup implements the IFrameProcessor interface.
func (p *AsyncFrameProcessor) Cleanup() {
	logger.Info("Cleanuping", "name", p.Name())
	p.interruptionMu.Lock()
	defer p.interruptionMu.Unlock()

	// Cancel the push frame task
	p.cancel()
	logger.Info("Cleanuping1", "name", p.Name())

	// Wait for the task to finish
	p.pushFrameTask.Wait()
	logger.Info("Cleanup Done", "name", p.Name())
}

// HandleInterruptions handles interruption frames.
func (p *AsyncFrameProcessor) HandleInterruptions(frame frames.Frame) {
	// out-of-band interruption handling
	//if !p.allowInterruptions {
	//	log.Printf("Warning: interruption frames are not allowed for processor %s", p.name)
	//	return
	//}

	p.interruptionMu.Lock()
	defer p.interruptionMu.Unlock()

	// Cancel the current task
	p.cancel()

	// Wait for the task to finish
	p.pushFrameTask.Wait()

	// Reset the task
	p.ctx, p.cancel = context.WithCancel(context.Background())

	// Push an out-of-band frame (not using the ordered push frame task)
	p.PushFrame(frame, FrameDirectionDownstream)

	// Create a new queue and task
	p.pushQueue = make(chan pushItem, p.pushQueueSize)
	p.pushFrameTask = &sync.WaitGroup{}
	p.createPushTask()
	logger.Info("AsyncFrameProcessor createPushTask is OK!")
}

// createPushTask creates a new push frame task.
func (p *AsyncFrameProcessor) createPushTask() {
	p.pushFrameTask.Add(1)
	go p.pushFrameTaskHandler()
}

// QueueFrame queues a frame for processing.
func (p *AsyncFrameProcessor) QueueFrame(frame frames.Frame, direction FrameDirection) {
	select {
	case p.pushQueue <- pushItem{frame: frame, direction: direction}:
	default:
		logger.Warnf("Warning: push queue is full for %s, frame: %+v direction: %s", p.name, frame, direction)
		time.Sleep(10 * time.Second)
	}
}
func (p *AsyncFrameProcessor) QueueUpStreamFrame(frame frames.Frame) {
	p.QueueFrame(frame, FrameDirectionUpstream)
}
func (p *AsyncFrameProcessor) QueueDownStreamFrame(frame frames.Frame) {
	p.QueueFrame(frame, FrameDirectionDownstream)
}

// pushFrameTaskHandler is the handler for the push frame task.
func (p *AsyncFrameProcessor) pushFrameTaskHandler() {
	defer p.pushFrameTask.Done()

	running := true
	for running {
		select {
		case <-p.ctx.Done():
			logger.Info(fmt.Sprintf("%s pushFrameTaskHandler cancelled", p.name))
			return
		case item, ok := <-p.pushQueue:
			if !ok {
				// Channel closed
				logger.Warn(fmt.Sprintf("%s push queue closed", p.name))
				return
			}

			// Push the frame
			p.PushFrame(item.frame, item.direction)

			// Check if this is an end frame
			if _, ok := item.frame.(frames.EndFrame); ok {
				running = false
			}
		case <-time.After(1 * time.Second):
			// Timeout, continue the loop
			continue
		}
	}
}
