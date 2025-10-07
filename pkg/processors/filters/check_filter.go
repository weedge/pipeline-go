package filters

import (
	"testing"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// CheckFilter is a processor for testing that asserts certain conditions.
type CheckFilter struct {
	processors.FrameProcessor
	t        *testing.T
	text     string
	checkImg bool
}

// NewCheckFilter creates a new CheckFilter.
func NewCheckFilter(t *testing.T, text string, checkImg bool) *CheckFilter {
	return &CheckFilter{
		t:        t,
		text:     text,
		checkImg: checkImg,
	}
}

// ProcessFrame asserts conditions on the frame.
func (c *CheckFilter) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	if textFrame, ok := frame.(*frames.TextFrame); ok {
		if c.text != "" {
			if textFrame.Text == c.text {
				c.t.Errorf("Text should not be: %s", c.text)
			}
		}
	}
	if c.checkImg {
		if _, ok := frame.(*frames.ImageRawFrame); ok {
			c.t.Error("Should not contain image frame")
		}
	}
	c.PushFrame(frame, direction)
}
