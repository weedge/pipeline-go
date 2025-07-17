package processors

import (
	"log"

	"github.com/weedge/pipeline-go/pkg/frames"
)

// PrintOutFrameProcessor prints frames and can send control frames upstream.
type PrintOutFrameProcessor struct {
	BaseProcessor
}

// NewPrintOutFrameProcessor creates a new PrintOutFrameProcessor.
func NewPrintOutFrameProcessor() *PrintOutFrameProcessor {
	return &PrintOutFrameProcessor{}
}

// ProcessFrame prints the frame and handles specific text commands.
func (p *PrintOutFrameProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	log.Printf("PrintOut: %s", frame)

	if textFrame, ok := frame.(*frames.TextFrame); ok {
		switch textFrame.Text {
		case "end":
			p.PushFrame(frames.NewEndFrame(), FrameDirectionUpstream)
			log.Println("End ok")
		case "cancel":
			p.PushFrame(frames.NewCancelFrame(), FrameDirectionUpstream)
			log.Println("Cancel ok")
		case "endTask":
			p.PushFrame(frames.NewStopTaskFrame(), FrameDirectionUpstream)
			log.Println("End Task ok")
		}
	}
	p.PushFrame(frame, direction)
}
