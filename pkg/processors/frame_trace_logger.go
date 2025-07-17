package processors

import (
	"log"

	"github.com/weedge/pipeline-go/pkg/frames"
)

type FrameTraceLogger struct {
	BaseProcessor
}

func NewFrameTraceLogger() *FrameTraceLogger {
	return &FrameTraceLogger{}
}

func (l *FrameTraceLogger) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	log.Printf("get %T: %+v", frame, frame)
	l.PushFrame(frame, direction)
}
