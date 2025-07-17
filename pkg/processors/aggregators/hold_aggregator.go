package aggregators

import (
	"reflect"
	"sync"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/notifiers"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// HoldLastFrameAggregator holds the last frame of a specific type until signaled to release it.
type HoldLastFrameAggregator struct {
	processors.BaseProcessor
	holdFrameType reflect.Type
	notifier      notifiers.Notifier
	lastFrame     frames.Frame
	lock          sync.Mutex
	once          sync.Once
}

// NewHoldLastFrameAggregator creates a new HoldLastFrameAggregator.
func NewHoldLastFrameAggregator(frameType reflect.Type, notifier notifiers.Notifier) *HoldLastFrameAggregator {
	return &HoldLastFrameAggregator{
		holdFrameType: frameType,
		notifier:      notifier,
	}
}

// startReleaseListener starts a goroutine that waits for a notification.
func (a *HoldLastFrameAggregator) startReleaseListener(direction processors.FrameDirection) {
	go func() {
		for range a.notifier.Wait() {
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
