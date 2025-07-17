package processors

import (
	"reflect"
	"sync"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
)

// IdleProcessor detects periods of inactivity and emits an IdleFrame.
type IdleProcessor struct {
	BaseProcessor
	timeout         time.Duration
	resetTypes      map[reflect.Type]bool
	timer           *time.Timer
	lock            sync.Mutex
	once            sync.Once
	downstreamDir   FrameDirection
}

// NewIdleProcessor creates a new IdleProcessor.
// If no resetTypes are provided, any frame will reset the timer.
func NewIdleProcessor(timeout time.Duration, resetTypes ...reflect.Type) *IdleProcessor {
	typeMap := make(map[reflect.Type]bool)
	for _, t := range resetTypes {
		typeMap[t] = true
	}
	return &IdleProcessor{
		timeout:    timeout,
		resetTypes: typeMap,
	}
}

func (p *IdleProcessor) startTimer() {
	p.timer = time.AfterFunc(p.timeout, func() {
		p.PushFrame(&frames.IdleFrame{}, p.downstreamDir)
		// After firing, the timer is stopped. We restart it to detect the next idle period.
		p.lock.Lock()
		p.timer.Reset(p.timeout)
		p.lock.Unlock()
	})
}

func (p *IdleProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// Start the timer on the first frame.
	p.once.Do(func() {
		p.downstreamDir = direction
		p.startTimer()
	})

	// Reset the timer if an appropriate frame arrives.
	shouldReset := len(p.resetTypes) == 0
	if !shouldReset {
		if _, found := p.resetTypes[reflect.TypeOf(frame)]; found {
			shouldReset = true
		}
	}

	if shouldReset {
		p.timer.Reset(p.timeout)
	}

	// Pass the original frame through.
	p.PushFrame(frame, direction)

	// Stop the timer on EndFrame.
	if _, ok := frame.(frames.EndFrame); ok {
		p.timer.Stop()
	}
}
