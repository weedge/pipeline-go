package processors

import (
	"fmt"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
)

// PrintOutFrameProcessor prints frames and can send control frames upstream.
type PrintOutFrameProcessor struct {
	FrameProcessor
}

// NewPrintOutFrameProcessor creates a new PrintOutFrameProcessor.
func NewPrintOutFrameProcessor() *PrintOutFrameProcessor {
	return &PrintOutFrameProcessor{}
}

// ProcessFrame prints the frame and handles specific text commands.
func (p *PrintOutFrameProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	logger.Info(fmt.Sprintf("PrintOut: %s", frame))

	if textFrame, ok := frame.(*frames.TextFrame); ok {
		switch textFrame.Text {
		case "end":
			p.PushFrame(frames.NewEndFrame(), FrameDirectionUpstream)
			logger.Info("End ok")
		case "cancel":
			p.PushFrame(frames.NewCancelFrame(), FrameDirectionUpstream)
			logger.Info("Cancel ok")
		case "endTask":
			p.PushFrame(frames.NewStopTaskFrame(), FrameDirectionUpstream)
			logger.Info("End Task ok")
		}
	}
	p.PushFrame(frame, direction)
}
