package pipeline

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/notifiers"
	"github.com/weedge/pipeline-go/pkg/processors"
	"github.com/weedge/pipeline-go/pkg/processors/aggregators"
	"github.com/weedge/pipeline-go/pkg/processors/filters"
)

func TestSimple(t *testing.T) {
	pipeline := NewPipeline([]processors.FrameProcessor{
		processors.NewFrameTraceLogger("simple", 0),
	}, nil, nil)
	fmt.Println(pipeline)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("你好"))
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{Width: 0, Height: 0}, "PNG", "RGB"))
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}

func TestParallelPipeline(t *testing.T) {
	pipeline := NewPipeline([]processors.FrameProcessor{
		processors.NewFrameTraceLogger("0", 0),
		NewParallelPipeline(
			[]processors.FrameProcessor{processors.NewFrameTraceLogger("1.0", 1000)},
			[]processors.FrameProcessor{processors.NewFrameTraceLogger("1.1", 0)},
		),
		processors.NewFrameTraceLogger("3", 0),
	}, nil, nil)
	fmt.Println(pipeline)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("你好"))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}

func TestSyncParallelPipeline(t *testing.T) {
	pipeline := NewPipeline([]processors.FrameProcessor{
		processors.NewFrameTraceLogger("0", 0),
		NewSyncParallelPipeline(
			NewPipeline([]processors.FrameProcessor{processors.NewFrameTraceLogger("1.0", 1000)}, nil, nil),
			NewPipeline([]processors.FrameProcessor{processors.NewFrameTraceLogger("1.1", 0)}, nil, nil),
		),
		processors.NewFrameTraceLogger("3", 0),
	}, nil, nil)
	fmt.Println(pipeline)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("你好"))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}

func TestFunctionFilter(t *testing.T) {
	text := "你好"

	textFilter := func(frame frames.Frame) bool {
		if textFrame, ok := frame.(*frames.TextFrame); ok {
			if textFrame.Text == text {
				return false
			}
		}
		return true
	}

	imageFilter := func(frame frames.Frame) bool {
		if _, ok := frame.(*frames.ImageRawFrame); ok {
			return false
		}
		return true
	}

	pipeline := NewPipeline([]processors.FrameProcessor{
		filters.NewFrameFilter(textFilter),
		filters.NewCheckFilter(t, text, false),
		filters.NewFrameFilter(imageFilter),
		filters.NewCheckFilter(t, text, true),
	}, nil, nil)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame(text))
	task.QueueFrame(frames.NewTextFrame("你好!"))
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{Width: 0, Height: 0}, "PNG", "RGB"))
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}

func TestTypeFilter(t *testing.T) {
	var collectedFrames []frames.Frame
	var mu sync.Mutex

	collector := processors.NewOutputProcessor(func(frame frames.Frame) {
		// We don't want to collect control frames
		switch frame.(type) {
		case *frames.StartFrame, *frames.EndFrame:
			return
		}
		mu.Lock()
		defer mu.Unlock()
		collectedFrames = append(collectedFrames, frame)
	})

	pipeline := NewPipeline([]processors.FrameProcessor{
		// This filter should only allow TextFrames and AudioRawFrames to pass.
		filters.NewTypeFilter([]interface{}{&frames.TextFrame{}, &frames.AudioRawFrame{}}),
		collector,
	}, nil, nil)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("one"))
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{Width: 0, Height: 0}, "PNG", "RGB"))
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))
	task.QueueFrame(frames.NewTextFrame("two"))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()

	// Assert that only the correct frames were collected
	assert.Len(t, collectedFrames, 3)
	assert.IsType(t, &frames.TextFrame{}, collectedFrames[0])
	assert.IsType(t, &frames.AudioRawFrame{}, collectedFrames[1])
	assert.IsType(t, &frames.TextFrame{}, collectedFrames[2])
}

func TestHoldFramesAggregator(t *testing.T) {
	notifier := notifiers.NewChannelNotifier()

	wakeNotifierFilter := func(frame frames.Frame) bool {
		if _, ok := frame.(*frames.SyncNotifyFrame); ok {
			go func() {
				time.Sleep(100 * time.Millisecond)
				notifier.Notify()
			}()
		}
		return true
	}

	aggregator := aggregators.NewHoldFramesAggregator(
		[]interface{}{&frames.TextFrame{}},
		notifier,
	)

	pipeline := NewPipeline([]processors.FrameProcessor{
		aggregator,
		filters.NewFrameFilter(wakeNotifierFilter),
		processors.NewPrintOutFrameProcessor(),
	}, nil, nil)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("Hello many frames, "))
	task.QueueFrame(frames.NewSyncNotifyFrame())
	task.QueueFrame(frames.NewImageRawFrame(
		[]byte{}, frames.ImageSize{Width: 0, Height: 0}, "JPEG", "RGB",
	))
	task.QueueFrame(frames.NewTextFrame("Goodbye1."))
	task.QueueFrame(frames.NewImageRawFrame(
		[]byte{}, frames.ImageSize{Width: 0, Height: 0}, "PNG", "RGB",
	))
	task.QueueFrame(frames.NewStopTaskFrame())

	task.Run()
}

func TestHoldLastFrameAggregator(t *testing.T) {
	notifier := notifiers.NewChannelNotifier()

	wakeNotifierFilter := func(frame frames.Frame) bool {
		if _, ok := frame.(*frames.SyncNotifyFrame); ok {
			go func() {
				time.Sleep(100 * time.Millisecond)
				notifier.Notify()
			}()
		}
		return true
	}

	aggregator := aggregators.NewHoldLastFrameAggregator(
		[]interface{}{&frames.TextFrame{}},
		notifier,
	)

	pipeline := NewPipeline([]processors.FrameProcessor{
		aggregator,
		filters.NewFrameFilter(wakeNotifierFilter),
		processors.NewPrintOutFrameProcessor(),
	}, nil, nil)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("Hello last one frame, "))
	task.QueueFrame(frames.NewSyncNotifyFrame())
	task.QueueFrame(frames.NewImageRawFrame(
		[]byte{}, frames.ImageSize{Width: 0, Height: 0}, "JPEG", "RGB",
	))
	task.QueueFrame(frames.NewTextFrame("Goodbye1."))
	task.QueueFrame(frames.NewImageRawFrame(
		[]byte{}, frames.ImageSize{Width: 0, Height: 0}, "PNG", "RGB",
	))
	task.QueueFrame(frames.NewTextFrame("end"))
	task.QueueFrame(frames.NewTextFrame("endTask"))
	task.QueueFrame(frames.NewStopTaskFrame())

	task.Run()
}

func TestGatedAggregator(t *testing.T) {
	var collectedFrames []frames.Frame
	var mu sync.Mutex

	collector := processors.NewOutputProcessor(func(frame frames.Frame) {
		// We don't want to collect control frames
		switch frame.(type) {
		case *frames.StartFrame, *frames.EndFrame:
			return
		}
		mu.Lock()
		defer mu.Unlock()
		collectedFrames = append(collectedFrames, frame)
	})

	gateOpenFn := func(f frames.Frame) bool {
		_, ok := f.(*frames.ImageRawFrame)
		return ok
	}
	gateCloseFn := func(f frames.Frame) bool {
		_, ok := f.(*frames.TextFrame)
		return ok
	}

	aggregator := aggregators.NewGatedAggregator(gateOpenFn, gateCloseFn, false, processors.FrameDirectionDownstream)

	pipeline := NewPipeline([]processors.FrameProcessor{
		aggregator,
		collector,
	}, nil, nil)

	task := NewPipelineTask(pipeline, PipelineParams{})

	// 1. Gate is closed. This frame is dropped.
	task.QueueFrame(frames.NewTextFrame("Hello, "))
	// 2. Gate is closed. This frame is dropped.
	task.QueueFrame(frames.NewTextFrame("Hello again."))
	// 3. Gate is closed, but this frame opens it. It is passed through.
	task.QueueFrame(frames.NewImageRawFrame(
		[]byte{}, frames.ImageSize{Width: 0, Height: 0}, "JPEG", "RGB",
	))
	// 4. Gate is open. This frame is passed through, and then closes the gate.
	task.QueueFrame(frames.NewTextFrame("Goodbye1."))
	// 5. Gate is closed. This frame is dropped.
	task.QueueFrame(frames.NewTextFrame("Goodbye2."))
	// 6. Gate is closed, but this frame opens it. It is passed through.
	task.QueueFrame(frames.NewImageRawFrame(
		[]byte{}, frames.ImageSize{Width: 0, Height: 0}, "JPEG", "RGB",
	))
	// 7. Gate is open. This frame is passed through, and then closes the gate.
	task.QueueFrame(frames.NewTextFrame("Goodbye3."))

	task.QueueFrame(frames.NewEndFrame())

	task.Run()

	// Assert that only the frames that passed through the open gate were collected
	assert.Len(t, collectedFrames, 4)
	assert.IsType(t, &frames.ImageRawFrame{}, collectedFrames[0])
	assert.IsType(t, &frames.TextFrame{}, collectedFrames[1])
	assert.IsType(t, &frames.ImageRawFrame{}, collectedFrames[2])
	assert.IsType(t, &frames.TextFrame{}, collectedFrames[3])
}