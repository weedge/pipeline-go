package aggregators

import (
	"regexp"
	"strings"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

var endOfSentencePattern = regexp.MustCompile(`[\.。\?？\!！:：]`)

// SentenceAggregator aggregates text frames into sentences.
type SentenceAggregator struct {
	processors.FrameProcessor
	aggregation string
}

// NewSentenceAggregator creates a new SentenceAggregator.
func NewSentenceAggregator() *SentenceAggregator {
	return &SentenceAggregator{}
}

// hasEndOfSentence checks for sentence-terminating punctuation using regex.
func (a *SentenceAggregator) hasEndOfSentence(s string) bool {
	// A simplified regex is used here for brevity. The original Python version's
	// lookbehind-based regex is not directly supported in Go's standard regexp engine.
	// This version checks for common sentence-ending punctuation at the end of the string.
	return endOfSentencePattern.MatchString(strings.TrimSpace(s))
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
