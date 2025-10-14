package processors

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
)

// LoggerProcessor is a simple processor that logs frames.
type LoggerProcessor struct {
	FrameProcessor
	name string
}

func NewLoggerProcessor(name string) *LoggerProcessor {
	return &LoggerProcessor{
		FrameProcessor: *NewFrameProcessor(name),
	}
}

func (p *LoggerProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	logger.Info(fmt.Sprintf("[%s] received frame: %+v, direction: %d", p.name, frame, direction))
	p.PushFrame(frame, direction)
}

// FrameLoggerProcessor 是一个更高级的日志处理器，可以过滤特定类型的帧
type FrameLoggerProcessor struct {
	FrameProcessor
	prefix            string
	ignoredFrameTypes []reflect.Type
	includeFrameTypes []reflect.Type
	maxIdToLogs       []uint64
}

// NewFrameLoggerProcessor 创建一个新的 FrameLoggerProcessor 实例
func NewFrameLoggerProcessor(
	name string,
	prefix string,
	ignoredFrameTypes []frames.Frame,
	includeFrameTypes []frames.Frame,
	maxIdToLogs []uint64,
) *FrameLoggerProcessor {
	ignoredTypes := make([]reflect.Type, len(ignoredFrameTypes))
	for i, frame := range ignoredFrameTypes {
		ignoredTypes[i] = reflect.TypeOf(frame)
	}

	includeTypes := make([]reflect.Type, len(includeFrameTypes))
	for i, frame := range includeFrameTypes {
		includeTypes[i] = reflect.TypeOf(frame)
	}

	return &FrameLoggerProcessor{
		FrameProcessor:    *NewFrameProcessor(name),
		prefix:            prefix,
		ignoredFrameTypes: ignoredTypes,
		includeFrameTypes: includeTypes,
		maxIdToLogs:       maxIdToLogs,
	}
}

// NewFrameLoggerProcessorWithMaxIdToLogs creates a new FrameLoggerProcessor with a maximum number of IDs to log
func NewFrameLoggerProcessorWithMaxIdToLogs(maxIdToLogs []uint64) *FrameLoggerProcessor {
	return NewFrameLoggerProcessor("FrameLoggerProcessor", "Frame", []frames.Frame{}, []frames.Frame{}, maxIdToLogs)
}

// NewDefaultFrameLoggerProcessor 创建一个带有默认设置的 FrameLoggerProcessor
func NewDefaultFrameLoggerProcessor() *FrameLoggerProcessor {
	return NewFrameLoggerProcessor("FrameLoggerProcessor", "Frame", []frames.Frame{}, []frames.Frame{}, []uint64{})
}

// NewDefaultFrameLoggerProcessorWithName 创建一个带有默认设置的 FrameLoggerProcessor
func NewDefaultFrameLoggerProcessorWithName(name string) *FrameLoggerProcessor {
	return NewFrameLoggerProcessor(name, "Frame", []frames.Frame{}, []frames.Frame{}, []uint64{})
}

// NewDefaultFrameLoggerProcessorWithIncludeFrame 创建一个包括include frames设置的 FrameLoggerProcessor
func NewDefaultFrameLoggerProcessorWithIncludeFrame(includeFrames []frames.Frame) *FrameLoggerProcessor {
	return NewFrameLoggerProcessor("FrameLoggerProcessor", "Frame", []frames.Frame{}, includeFrames, []uint64{})
}

// NewDefaultFrameLoggerProcessorWithIngoreFrame 创建一个忽略ignore frames设置的 FrameLoggerProcessor
func NewDefaultFrameLoggerProcessorWithIngoreFrame(ingoreFrames []frames.Frame) *FrameLoggerProcessor {
	return NewFrameLoggerProcessor("FrameLoggerProcessor", "Frame", ingoreFrames, []frames.Frame{}, []uint64{})
}

func (p *FrameLoggerProcessor) WithMaxIdToLogs(maxIdToLogs []uint64) *FrameLoggerProcessor {
	p.maxIdToLogs = maxIdToLogs
	return p
}

func (p *FrameLoggerProcessor) WithIncludeFrame(includeFrameTypes []frames.Frame) *FrameLoggerProcessor {
	includeTypes := make([]reflect.Type, len(includeFrameTypes))
	for i, frame := range includeFrameTypes {
		includeTypes[i] = reflect.TypeOf(frame)
	}
	p.includeFrameTypes = includeTypes
	return p
}

func (p *FrameLoggerProcessor) WithIgnoreFrame(ignoreFrameTypes []frames.Frame) *FrameLoggerProcessor {
	ignoreTypes := make([]reflect.Type, len(ignoreFrameTypes))
	for i, frame := range ignoreFrameTypes {
		ignoreTypes[i] = reflect.TypeOf(frame)
	}
	p.ignoredFrameTypes = ignoreTypes
	return p
}

func (p *FrameLoggerProcessor) WithPrefix(prefix string) *FrameLoggerProcessor {
	p.prefix = prefix
	return p
}

func (p *FrameLoggerProcessor) Prefix() string {
	return p.prefix
}

// isIgnoredFrame 检查帧是否应该被忽略
func (p *FrameLoggerProcessor) isIgnoredFrame(frame frames.Frame) bool {
	if len(p.ignoredFrameTypes) == 0 {
		return false
	}

	frameType := reflect.TypeOf(frame)
	return slices.Contains(p.ignoredFrameTypes, frameType)
}

// isIncludedFrame 检查帧是否在包含列表中
func (p *FrameLoggerProcessor) isIncludedFrame(frame frames.Frame) bool {
	return p.indexIncludedFrame(frame) >= 0
}

// isIncludedFrame 检查帧是否在包含列表中
func (p *FrameLoggerProcessor) indexIncludedFrame(frame frames.Frame) int {
	if len(p.includeFrameTypes) == 0 {
		return -1
	}

	frameType := reflect.TypeOf(frame)
	return slices.Index(p.includeFrameTypes, frameType)
}

func (p *FrameLoggerProcessor) ProcessFrame(frame frames.Frame, direction FrameDirection) {
	// 检查是否应该处理此帧
	if !p.isIgnoredFrame(frame) {
		index := p.indexIncludedFrame(frame)
		if index >= 0 {
			fromTo := p.name

			// 只有当 prev 和 next 不为 nil 时才使用它们
			if p.prev != nil && p.next != nil {
				switch direction {
				case FrameDirectionDownstream:
					fromTo = fmt.Sprintf("%s(%T) ---> %s(%T)", p.prev.Name(), p.prev, p.name, p)
				case FrameDirectionUpstream:
					fromTo = fmt.Sprintf("%s(%T) <--- %s(%T)", p.name, p, p.next.Name(), p.next)
				}
			} else if p.prev != nil {
				fromTo = fmt.Sprintf("%s(%T) ---> %s(%T)", p.prev.Name(), p.prev, p.name, p)
			} else if p.next != nil {
				fromTo = fmt.Sprintf("%s(%T) ---> %s(%T)", p.name, p, p.next.Name(), p.next)
			}

			msg := fmt.Sprintf("%s %s: (%T)%s", fromTo, p.prefix, frame, frame.String())

			if len(p.maxIdToLogs) == 0 || (len(p.maxIdToLogs) > index && frame.ID()%p.maxIdToLogs[index] == 0) {
				logger.Info(msg)
			}
		}
	}

	p.PushFrame(frame, direction)
}
