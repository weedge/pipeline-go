package aggregators

import (
	"reflect"
	"strings"

	"github.com/weedge/pipeline-go/pkg"
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// SentenceAggregator aggregates text frames into sentences.
type SentenceAggregator struct {
	*processors.FrameProcessor
	aggregation string
	endFrame    reflect.Type // endFrame to flush sentence
}

// NewSentenceAggregator creates a new SentenceAggregator.
func NewSentenceAggregator() *SentenceAggregator {
	return NewSentenceAggregatorWithEnd(reflect.TypeOf(&frames.EndFrame{}))
}

func NewSentenceAggregatorWithEnd(endFrame reflect.Type) *SentenceAggregator {
	return &SentenceAggregator{
		FrameProcessor: processors.NewFrameProcessor("SentenceAggregator"),
		endFrame:       endFrame,
	}
}

// hasEndOfSentence checks for sentence-terminating punctuation using regex.
func (a *SentenceAggregator) hasEndOfSentence(s string) bool {
	return pkg.MatchEndOfSentence(strings.TrimSpace(s))
}

// ProcessFrame accumulates text and emits a frame when a sentence is complete.
func (a *SentenceAggregator) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	isPushAgg := reflect.TypeOf(frame) == a.endFrame

	switch f := frame.(type) {
	case *frames.TextFrame:
		// Add a space if the aggregation is not empty and the new text doesn't start with one.
		a.aggregation += f.Text
		if a.hasEndOfSentence(a.aggregation) {
			a.PushFrame(frames.NewTextFrame(a.aggregation), direction)
			a.aggregation = ""
		}
	case *frames.EndFrame:
		// If there's any leftover text, push it out before ending.
		if a.aggregation != "" {
			a.PushFrame(frames.NewTextFrame(a.aggregation), direction)
			a.aggregation = ""
		}
		a.PushFrame(f, direction)
	default:
		// Pass other frame types through.
		a.PushFrame(frame, direction)
	}

	if isPushAgg && a.aggregation != "" {
		a.PushFrame(frames.NewTextFrame(a.aggregation), direction)
		a.aggregation = ""
	}
}
