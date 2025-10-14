package aggregators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// mockProcessor is a simple processor for testing.
type mockProcessor struct {
	*processors.FrameProcessor
	receivedFrames map[processors.FrameDirection][]frames.Frame
}

func NewMockProcessor() *mockProcessor {
	return &mockProcessor{
		FrameProcessor: processors.NewFrameProcessor("mock_processor"),
		receivedFrames: make(map[processors.FrameDirection][]frames.Frame),
	}
}
func (p *mockProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	p.receivedFrames[direction] = append(p.receivedFrames[direction], frame)
	// Push frame further in mock processor to simulate a real pipeline
	// But only if there is a next processor
	p.FrameProcessor.PushFrame(frame, direction)
}

func (p *mockProcessor) GetReceivedFrames(direction processors.FrameDirection) []frames.Frame {
	return p.receivedFrames[direction]
}

func TestSentenceAggregator_ProcessFrame_TextFrames(t *testing.T) {
	// Create a mock processor to receive frames
	mockProc := NewMockProcessor()

	// Create a sentence aggregator and link it to the mock processor
	sentenceAggregator := NewSentenceAggregator()
	sentenceAggregator.Link(mockProc)

	// Send text frames that form a complete sentence
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "Hello"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "world"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "."}, processors.FrameDirectionDownstream)

	// Check that the aggregated sentence was sent
	received := mockProc.GetReceivedFrames(processors.FrameDirectionDownstream)
	assert.Equal(t, 1, len(received))
	textFrame, ok := received[0].(*frames.TextFrame)
	assert.True(t, ok)
	assert.Equal(t, "Hello world .", textFrame.Text)

	// Check that aggregation was reset
	assert.Equal(t, "", sentenceAggregator.aggregation)
}

func TestSentenceAggregator_ProcessFrame_TextFramesWithSpaces(t *testing.T) {
	// Create a mock processor to receive frames
	mockProc := NewMockProcessor()

	// Create a sentence aggregator and link it to the mock processor
	sentenceAggregator := NewSentenceAggregator()
	sentenceAggregator.Link(mockProc)

	// Send text frames with spaces
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "Hello"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: " beautiful"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: " world"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "!"}, processors.FrameDirectionDownstream)

	// Check that the aggregated sentence was sent
	received := mockProc.GetReceivedFrames(processors.FrameDirectionDownstream)
	assert.Equal(t, 1, len(received))
	textFrame, ok := received[0].(*frames.TextFrame)
	assert.True(t, ok)
	assert.Equal(t, "Hello beautiful world !", textFrame.Text)
}

func TestSentenceAggregator_ProcessFrame_MultipleSentences(t *testing.T) {
	// Create a mock processor to receive frames
	mockProc := NewMockProcessor()

	// Create a sentence aggregator and link it to the mock processor
	sentenceAggregator := NewSentenceAggregator()
	sentenceAggregator.Link(mockProc)

	// Send text frames that form multiple sentences
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "Hello"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "world"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "."}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "How"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "are"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "you"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "?"}, processors.FrameDirectionDownstream)

	// Check that both sentences were sent
	received := mockProc.GetReceivedFrames(processors.FrameDirectionDownstream)
	assert.Equal(t, 2, len(received))

	textFrame1, ok := received[0].(*frames.TextFrame)
	assert.True(t, ok)
	assert.Equal(t, "Hello world .", textFrame1.Text)

	textFrame2, ok := received[1].(*frames.TextFrame)
	assert.True(t, ok)
	assert.Equal(t, "How are you ?", textFrame2.Text)

	// Check that aggregation was reset
	assert.Equal(t, "", sentenceAggregator.aggregation)
}

func TestSentenceAggregator_ProcessFrame_EndFrameWithLeftoverText(t *testing.T) {
	// Create a mock processor to receive frames
	mockProc := NewMockProcessor()

	// Create a sentence aggregator and link it to the mock processor
	sentenceAggregator := NewSentenceAggregator()
	sentenceAggregator.Link(mockProc)

	// Send text frames without ending punctuation
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "Hello"}, processors.FrameDirectionDownstream)
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "world"}, processors.FrameDirectionDownstream)

	// Send an end frame
	endFrame := frames.NewEndFrame()
	sentenceAggregator.ProcessFrame(endFrame, processors.FrameDirectionDownstream)

	// Check that the leftover text was sent followed by the end frame
	received := mockProc.GetReceivedFrames(processors.FrameDirectionDownstream)

	assert.Equal(t, 2, len(received))

	// The first frame should be the aggregated text
	textFrame, ok := received[0].(*frames.TextFrame)
	assert.True(t, ok)
	assert.Equal(t, "Hello world", textFrame.Text)

	// The second frame should be the end frame
	_, ok = received[1].(*frames.EndFrame)
	assert.True(t, ok)

	// Check that aggregation was reset
	assert.Equal(t, "", sentenceAggregator.aggregation)
}

func TestSentenceAggregator_ProcessFrame_PassThroughOtherFrames(t *testing.T) {
	// Create a mock processor to receive frames
	mockProc := NewMockProcessor()

	// Create a sentence aggregator and link it to the mock processor
	sentenceAggregator := NewSentenceAggregator()
	sentenceAggregator.Link(mockProc)

	// Send a start frame (should pass through)
	startFrame := frames.NewStartFrame()
	sentenceAggregator.ProcessFrame(startFrame, processors.FrameDirectionDownstream)

	// Send a text frame
	sentenceAggregator.ProcessFrame(&frames.TextFrame{Text: "Hello."}, processors.FrameDirectionDownstream)

	// Check that both frames were sent
	received := mockProc.GetReceivedFrames(processors.FrameDirectionDownstream)
	assert.Equal(t, 2, len(received))

	_, ok := received[0].(*frames.StartFrame)
	assert.True(t, ok)

	textFrame, ok := received[1].(*frames.TextFrame)
	assert.True(t, ok)
	assert.Equal(t, "Hello.", textFrame.Text)
}

func TestSentenceAggregator_hasEndOfSentence(t *testing.T) {
	aggregator := NewSentenceAggregator()

	// Test cases with ending punctuation
	assert.True(t, aggregator.hasEndOfSentence("."))
	assert.True(t, aggregator.hasEndOfSentence("。"))
	assert.True(t, aggregator.hasEndOfSentence("?"))
	assert.True(t, aggregator.hasEndOfSentence("？"))
	assert.True(t, aggregator.hasEndOfSentence("!"))
	assert.True(t, aggregator.hasEndOfSentence("！"))
	assert.True(t, aggregator.hasEndOfSentence(":"))
	assert.True(t, aggregator.hasEndOfSentence("："))

	// Test cases with whitespace
	assert.True(t, aggregator.hasEndOfSentence(" . "))
	assert.True(t, aggregator.hasEndOfSentence(" ? "))
	assert.True(t, aggregator.hasEndOfSentence(" ! "))

	// Test cases without ending punctuation
	assert.False(t, aggregator.hasEndOfSentence("hello"))
	assert.False(t, aggregator.hasEndOfSentence("hello world"))
	assert.False(t, aggregator.hasEndOfSentence(""))
	assert.False(t, aggregator.hasEndOfSentence("   "))
}
