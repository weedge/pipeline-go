package processors

import (
	"log"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
)

type FrameTraceLogger struct {
	BaseProcessor
	tag     string
	delayMs int
}

func NewFrameTraceLogger(tag string, delayMs int) *FrameTraceLogger {
	return &FrameTraceLogger{tag: tag, delayMs: delayMs}
}

func (l *FrameTraceLogger) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	if l.delayMs > 0 {
		time.Sleep(time.Duration(l.delayMs) * time.Millisecond)
	}
	log.Printf("Tag: %s; get Frame: %s", l.tag, frame)
	l.PushFrame(frame, direction)
}