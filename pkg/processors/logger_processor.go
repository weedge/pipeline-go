package processors

import (
	"fmt"
	"log/slog"

	"github.com/weedge/pipeline-go/pkg/frames"
)

// LoggerProcessor is a simple processor that logs frames.
type LoggerProcessor struct {
	FrameProcessor
	name string
}

func NewLoggerProcessor(name string) *LoggerProcessor {
	return &LoggerProcessor{name: name}
}

func (p *LoggerProcessor) Name() string {
	return p.name
}

func (p *LoggerProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	slog.Info(fmt.Sprintf("[%s] received frame: %+v, direction: %d", p.name, frame, direction))
	p.PushFrame(frame, direction)
}
