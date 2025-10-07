package processors

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weedge/pipeline-go/pkg/frames"
)

// mockProcessor is a simple processor for testing.
type mockProcessor struct {
	*FrameProcessor
	receivedFrames map[FrameDirection][]frames.Frame
}

func NewMockProcessor() *mockProcessor {
	return &mockProcessor{
		FrameProcessor: NewFrameProcessor("mock_processor"),
		receivedFrames: make(map[FrameDirection][]frames.Frame),
	}
}
func NewMockProcessorWithName(name string) *mockProcessor {
	return &mockProcessor{
		FrameProcessor: NewFrameProcessor(name),
		receivedFrames: make(map[FrameDirection][]frames.Frame),
	}
}
func (p *mockProcessor) Link(next IFrameProcessor) {
	p.next = next
	if p.verbose {
		log.Printf("%s(%T) -> %s(%T)", p.Name(), p, next.Name(), next)
	}
}

// SetPrev sets the previous processor in the pipeline.
func (p *mockProcessor) SetPrev(prev IFrameProcessor) {
	if p.verbose {
		log.Printf("%s(%T) <- %s(%T)", prev.Name(), prev, p.Name(), p)
	}
	p.prev = prev
}

func (p *mockProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	p.receivedFrames[direction] = append(p.receivedFrames[direction], frame)
	println("Added frame to receivedFrames, direction:", direction, "count:", len(p.receivedFrames[direction]), "processor:", p.Name())
	p.PushFrame(frame, direction)
}

func (p *mockProcessor) GetReceivedFrames() map[FrameDirection][]frames.Frame {
	return p.receivedFrames
}

func (p *mockProcessor) GetReceivedDirectionDownstreamFrames() []frames.Frame {
	return p.receivedFrames[FrameDirectionDownstream]
}
func (p *mockProcessor) GetReceivedDirectionUpstreamFrames() []frames.Frame {
	return p.receivedFrames[FrameDirectionUpstream]
}

func TestAsyncFrameProcessor(t *testing.T) {
	// asyncProc -> mockProc <-> mockProc1
	// Create a mock processor to receive frames
	mockProc := NewMockProcessor()
	mockProc.SetVerbose(true)
	mockProc1 := NewMockProcessorWithName("mock_processor1")
	mockProc1.SetVerbose(true)
	mockProc.Link(mockProc1)
	mockProc1.SetPrev(mockProc)

	// Create an async processor and link it to the mock processor
	asyncProc := NewAsyncFrameProcessor("async_processor")
	asyncProc.Link(mockProc)

	// Send some frames
	s_frame := frames.NewStartFrame()
	frame1 := frames.NewBaseFrameWithName("frame1")
	frame2 := frames.NewBaseFrameWithName("frame2")
	frame3 := frames.NewBaseFrameWithName("frame3")
	e_frame := frames.NewEndFrame()

	println("Sending frames...")
	asyncProc.ProcessFrame(s_frame, FrameDirectionDownstream)
	asyncProc.ProcessFrame(frame1, FrameDirectionDownstream)
	// Use PushFrame instead of ProcessFrame to ensure proper upstream handling
	println("Pushing upstream frame2")
	mockProc1.PushFrame(frame2, FrameDirectionUpstream)
	println("Pushing upstream frame3")
	mockProc1.PushFrame(frame3, FrameDirectionUpstream)
	asyncProc.ProcessFrame(e_frame, FrameDirectionDownstream)

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	// Cleanup
	asyncProc.Cleanup()

	// Check that frames were received by mockProc (middle of the chain)
	received_mock := mockProc.GetReceivedDirectionDownstreamFrames()
	println("mockProc downstream frames:", len(received_mock))

	up_received_mock := mockProc.GetReceivedDirectionUpstreamFrames()
	println("mockProc upstream frames:", len(up_received_mock))

	// Print all received frames for debugging
	println("All received frames in mockProc:")
	for direction, frames := range mockProc.receivedFrames {
		println("  Direction", direction, "count:", len(frames))
		for i, frame := range frames {
			println("    ", i, frame.Name())
		}
	}

	assert.Equal(t, 3, len(received_mock))
	assert.IsType(t, &frames.StartFrame{}, received_mock[0])
	assert.IsType(t, &frames.BaseFrame{}, received_mock[1])
	assert.IsType(t, &frames.EndFrame{}, received_mock[2])

	assert.Equal(t, 2, len(up_received_mock))
	assert.IsType(t, &frames.BaseFrame{}, up_received_mock[0])
	assert.IsType(t, &frames.BaseFrame{}, up_received_mock[1])

	// Check that frames were received by mockProc1 (end of the chain)
	received_mock1 := mockProc1.GetReceivedDirectionDownstreamFrames()
	println("mockProc1 downstream frames:", len(received_mock1))

	up_received_mock1 := mockProc1.GetReceivedDirectionUpstreamFrames()
	println("mockProc1 upstream frames:", len(up_received_mock1))

	// Print all received frames for debugging
	println("All received frames in mockProc1:")
	for direction, frames := range mockProc1.receivedFrames {
		println("  Direction", direction, "count:", len(frames))
		for i, frame := range frames {
			println("    ", i, frame.Name())
		}
	}

	assert.Equal(t, 3, len(received_mock1))
	assert.IsType(t, &frames.StartFrame{}, received_mock1[0])
	assert.IsType(t, &frames.BaseFrame{}, received_mock1[1])
	assert.IsType(t, &frames.EndFrame{}, received_mock1[2])

	assert.Equal(t, 0, len(up_received_mock1))
}

func TestInterruptionAsyncFrameProcessor(t *testing.T) {
	// asyncProc -> mockProc

	// Create a mock processor to receive frames
	mockProc := NewMockProcessor()

	// Create an async processor and link it to the mock processor
	asyncProc := NewAsyncFrameProcessor("interrupt_processor")
	asyncProc.SetVerbose(true)
	asyncProc.Link(mockProc)

	// Send an interruption frame
	asyncProc.ProcessFrame(frames.NewStartInterruptionFrame(), FrameDirectionDownstream)

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	// Send another frame after interruption
	frame := frames.NewBaseFrame()
	asyncProc.ProcessFrame(frame, FrameDirectionDownstream)

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	// Cleanup
	asyncProc.Cleanup()

	// Check that frames were received
	received := mockProc.GetReceivedFrames()
	for k, fs := range received {
		for _, frame := range fs {
			println(k.String(), frame.String())
		}
	}
	received_down := mockProc.GetReceivedDirectionDownstreamFrames()
	println(len(received_down))
	assert.GreaterOrEqual(t, len(received_down), 1) // At least the interruption frame should be received
}
