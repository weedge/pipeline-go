package processors

import (
	"fmt"
	"log/slog"

	"github.com/weedge/pipeline-go/pkg/frames"
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
	slog.Info(fmt.Sprintf("PrintOut: %s", frame))

	if textFrame, ok := frame.(*frames.TextFrame); ok {
		switch textFrame.Text {
		case "end":
			p.PushFrame(frames.NewEndFrame(), FrameDirectionUpstream)
			slog.Info("End ok")
		case "cancel":
			p.PushFrame(frames.NewCancelFrame(), FrameDirectionUpstream)
			slog.Info("Cancel ok")
		case "endTask":
			p.PushFrame(frames.NewStopTaskFrame(), FrameDirectionUpstream)
			slog.Info("End Task ok")
		}
	}
	p.PushFrame(frame, direction)
}

// getActualProcessor returns the actual processor instance.
func (p *PrintOutFrameProcessor) getActualProcessor() IFrameProcessor {
	return p
}
