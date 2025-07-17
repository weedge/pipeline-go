package processors

import (
	"log"
	"github.com/wuyong/pipeline-go/pkg/frames"
)

// LoggerProcessor is a simple processor that logs frames.
type LoggerProcessor struct {
	BaseProcessor
	Name string
}

func NewLoggerProcessor(name string) *LoggerProcessor {
	return &LoggerProcessor{Name: name}
}

func (p *LoggerProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	log.Printf("[%s] received frame: %+v, direction: %d", p.Name, frame, direction)
	p.PushFrame(frame, direction)
}
