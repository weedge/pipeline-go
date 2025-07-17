package aggregators

import (
	"reflect"
	"sync"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

// HoldLastFrameAggregator holds the last frame of a specific type until signaled to release it.
type HoldLastFrameAggregator struct {
	processors.BaseProcessor
	holdFrameType reflect.Type
	lastFrame     frames.Frame
	lock          sync.Mutex
	releaseChan   chan struct{}
	once          sync.Once
}

// NewHoldLastFrameAggregator creates a new HoldLastFrameAggregator.
// The frameType argument should be a reflect.Type of the frame to hold (e.g., reflect.TypeOf(&frames.TextFrame{})).
func NewHoldLastFrameAggregator(frameType reflect.Type) *HoldLastFrameAggregator {
	return &HoldLastFrameAggregator{
		holdFrameType: frameType,
		releaseChan:   make(chan struct{}),
	}
}

// Release signals the aggregator to release the held frame.
func (a *HoldLastFrameAggregator) Release() {
	a.releaseChan <- struct{}{}
}

// startReleaseListener starts a goroutine that waits for a release signal.
func (a *HoldLastFrameAggregator) startReleaseListener(direction processors.FrameDirection) {
	go func() {
		for range a.releaseChan {
			a.lock.Lock()
			if a.lastFrame != nil {
				a.PushFrame(a.lastFrame, direction)
				a.lastFrame = nil
			}
			a.lock.Unlock()
		}
	}()
}

func (a *HoldLastFrameAggregator) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Start the release listener only once.
	a.once.Do(func() {
		a.startReleaseListener(direction)
	})

	// Check if the incoming frame is of the type we want to hold.
	if reflect.TypeOf(frame) == a.holdFrameType {
		a.lock.Lock()
		a.lastFrame = frame
		a.lock.Unlock()
	} else {
		// Pass all other frames through immediately.
		a.PushFrame(frame, direction)
	}
}
