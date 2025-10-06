package aggregators

import (
	"reflect"
	"sync"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/notifiers"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// HoldFramesAggregator holds all frames of a specific type until signaled to release them.
type HoldFramesAggregator struct {
	processors.FrameProcessor
	holdFrameTypes []reflect.Type
	notifier       notifiers.Notifier
	heldFrames     []frames.Frame
	lock           sync.Mutex
	once           sync.Once
}

// NewHoldFramesAggregator creates a new HoldFramesAggregator.
func NewHoldFramesAggregator(frameTypes []interface{}, notifier notifiers.Notifier) *HoldFramesAggregator {
	var types []reflect.Type
	for _, t := range frameTypes {
		types = append(types, reflect.TypeOf(t))
	}
	return &HoldFramesAggregator{
		holdFrameTypes: types,
		notifier:       notifier,
	}
}

func (a *HoldFramesAggregator) startReleaseListener(direction processors.FrameDirection) {
	go func() {
		for range a.notifier.Wait() {
			a.lock.Lock()
			for _, frame := range a.heldFrames {
				a.PushFrame(frame, direction)
			}
			a.heldFrames = nil // Clear after releasing
			a.lock.Unlock()
		}
	}()
}

func (a *HoldFramesAggregator) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	a.once.Do(func() {
		a.startReleaseListener(direction)
	})

	isHeldType := false
	frameType := reflect.TypeOf(frame)
	for _, t := range a.holdFrameTypes {
		if frameType == t {
			isHeldType = true
			break
		}
	}

	if isHeldType {
		a.lock.Lock()
		a.heldFrames = append(a.heldFrames, frame)
		a.lock.Unlock()
	} else {
		a.PushFrame(frame, direction)
	}
}

// HoldLastFrameAggregator holds the last frame of a specific type until signaled to release it.
type HoldLastFrameAggregator struct {
	processors.FrameProcessor
	holdFrameTypes []reflect.Type
	notifier       notifiers.Notifier
	lastFrame      frames.Frame
	lock           sync.Mutex
	once           sync.Once
}

// NewHoldLastFrameAggregator creates a new HoldLastFrameAggregator.
func NewHoldLastFrameAggregator(frameTypes []interface{}, notifier notifiers.Notifier) *HoldLastFrameAggregator {
	var types []reflect.Type
	for _, t := range frameTypes {
		types = append(types, reflect.TypeOf(t))
	}
	return &HoldLastFrameAggregator{
		holdFrameTypes: types,
		notifier:       notifier,
	}
}

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
	a.once.Do(func() {
		a.startReleaseListener(direction)
	})

	isHeldType := false
	frameType := reflect.TypeOf(frame)
	for _, t := range a.holdFrameTypes {
		if frameType == t {
			isHeldType = true
			break
		}
	}

	if isHeldType {
		a.lock.Lock()
		a.lastFrame = frame
		a.lock.Unlock()
	} else {
		a.PushFrame(frame, direction)
	}
}
