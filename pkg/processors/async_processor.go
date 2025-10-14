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
	pushUpQueueSize       int
	pushUpQueue           chan pushItem
	pushUpFrameTask       *sync.WaitGroup
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
	pushQueueSize, pushUpQueueSize := 128, 0
	return NewAsyncFrameProcessorWithPushQueueSize(name, pushQueueSize, pushUpQueueSize)
}

func NewAsyncFrameProcessorWithPushQueueSize(name string, pushQueueSize, pushUpQueueSize int) *AsyncFrameProcessor {
	if pushQueueSize <= 0 {
		pushQueueSize = 128
	}
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

	if pushUpQueueSize > 0 {
		p.pushUpQueueSize = pushUpQueueSize
		p.pushUpQueue = make(chan pushItem, pushUpQueueSize)
		p.pushUpFrameTask = &sync.WaitGroup{}
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

	// Wait for the task to finish
	p.pushFrameTask.Wait()
	if p.pushUpQueueSize > 0 {
		p.pushUpFrameTask.Wait()
	}

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
	if p.pushUpQueueSize > 0 {
		p.pushUpFrameTask.Wait()
	}

	// Reset the task
	p.ctx, p.cancel = context.WithCancel(context.Background())

	// Push an out-of-band frame (not using the ordered push frame task)
	p.PushFrame(frame, FrameDirectionDownstream)

	// Create a new queue and task
	p.pushQueue = make(chan pushItem, p.pushQueueSize)
	p.pushFrameTask = &sync.WaitGroup{}
	p.createPushTask()
	logger.Info("AsyncFrameProcessor createPushTask is OK!")

	if p.pushUpQueueSize > 0 {
		// Create a new upstream queue and task
		p.pushUpQueue = make(chan pushItem, p.pushUpQueueSize)
		p.pushUpFrameTask = &sync.WaitGroup{}
		p.createPushTask()
		logger.Info("AsyncFrameProcessor createPushTask is OK!")
	}
}

// createPushTask creates a new push frame task.
func (p *AsyncFrameProcessor) createPushTask() {
	if p.pushFrameTask != nil {
		p.pushFrameTask.Add(1)
		go p.pushFrameTaskHandler()
		logger.Infof("%s create pushFrameTaskHandler", p.Name())
	}

	if p.pushUpFrameTask != nil {
		p.pushUpFrameTask.Add(1)
		go p.pushUpFrameTaskHandler()
		logger.Infof("%s create pushUpFrameTaskHandler", p.Name())
	}
}

// QueueFrame queues a frame for processing.
func (p *AsyncFrameProcessor) QueueFrame(frame frames.Frame, direction FrameDirection) {
	queue := p.pushQueue
	if p.pushUpQueueSize > 0 && direction == FrameDirectionUpstream {
		queue = p.pushUpQueue
	}

	select {
	case queue <- pushItem{frame: frame, direction: direction}:
	default:
		if p.pushUpQueueSize > 0 && direction == FrameDirectionUpstream {
			logger.Warnf("Warning: pushUpQueue is full for %s, frame: %+v direction: %s", p.name, frame, direction)
		} else {
			logger.Warnf("Warning: pushQueue is full for %s, frame: %+v direction: %s", p.name, frame, direction)
		}
		//time.Sleep(3 * time.Second)// open to test slow process
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
			// !NOTE:
			// - if PushFrame is Slow(e.g.: local llm gen token slow),
			// 	pushQueue maybe is full when push BotSpeakingFrame to upstream
			// - use param pushUpQueueSize>0 to create upstream task queue
			// !TODONE:
			// so need have two task queue: upstreamTask and downstreamTask,
			// But need think one Prcessor to do Up/Down Frame at the same time
			// if no shared stat, is ok.
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

// pushUpFrameTaskHandler is the handler for the push upstream frame task.
func (p *AsyncFrameProcessor) pushUpFrameTaskHandler() {
	defer p.pushUpFrameTask.Done()

	running := true
	for running {
		select {
		case <-p.ctx.Done():
			logger.Info(fmt.Sprintf("%s pushUpFrameTaskHandler cancelled", p.name))
			return
		case item, ok := <-p.pushUpQueue:
			if !ok {
				// Channel closed
				logger.Warn(fmt.Sprintf("%s push UpQueue closed", p.name))
				return
			}

			p.PushUpstreamFrame(item.frame)

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
