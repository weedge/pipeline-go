package processors

import (
	"fmt"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
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
	logger.Info(fmt.Sprintf("[%s] received frame: %+v, direction: %d", p.name, frame, direction))
	p.PushFrame(frame, direction)
}
