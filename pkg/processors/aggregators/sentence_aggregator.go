package aggregators

import (
	"strings"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

// SentenceAggregator aggregates text frames into sentences.
type SentenceAggregator struct {
	processors.BaseProcessor
	aggregation string
}

// NewSentenceAggregator creates a new SentenceAggregator.
func NewSentenceAggregator() *SentenceAggregator {
	return &SentenceAggregator{}
}

// hasEndOfSentence checks for sentence-terminating punctuation.
func (a *SentenceAggregator) hasEndOfSentence(s string) bool {
	return strings.HasSuffix(s, ".") || strings.HasSuffix(s, "?") || strings.HasSuffix(s, "!")
}

// ProcessFrame accumulates text and emits a frame when a sentence is complete.
func (a *SentenceAggregator) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	switch f := frame.(type) {
	case *frames.TextFrame:
		// Add a space if the aggregation is not empty and the new text doesn't start with one.
		if a.aggregation != "" && !strings.HasPrefix(f.Text, " ") {
			a.aggregation += " "
		}
		a.aggregation += f.Text
		if a.hasEndOfSentence(a.aggregation) {
			a.PushFrame(&frames.TextFrame{Text: a.aggregation}, direction)
			a.aggregation = ""
		}
	case frames.EndFrame:
		// If there's any leftover text, push it out before ending.
		if a.aggregation != "" {
			a.PushFrame(&frames.TextFrame{Text: a.aggregation}, direction)
			a.aggregation = ""
		}
		a.PushFrame(f, direction)
	default:
		// Pass other frame types through.
		a.PushFrame(frame, direction)
	}
}
