package processors

import (
	"github.com/weedge/pipeline-go/pkg/frames"
)

// StatelessTextTransformer applies a function to the text of a TextFrame.
type StatelessTextTransformer struct {
	FrameProcessor
	transformFn func(string) string
}

// NewStatelessTextTransformer creates a new StatelessTextTransformer.
func NewStatelessTextTransformer(transformFn func(string) string) *StatelessTextTransformer {
	return &StatelessTextTransformer{
		transformFn: transformFn,
	}
}

// ProcessFrame applies the transform function if the frame is a TextFrame.
func (t *StatelessTextTransformer) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	if textFrame, ok := frame.(*frames.TextFrame); ok {
		textFrame.Text = t.transformFn(textFrame.Text)
		t.PushFrame(textFrame, direction)
	} else {
		t.PushFrame(frame, direction)
	}
}
