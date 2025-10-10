package main

import (
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func main() {
	// 初始化日志配置
	logger.InitLoggerWithConfig(logger.NewDefaultLoggerConfig())

	// 0. 创建一个简单的帧
	frame := frames.NewBaseFrameWithName("TestFrame")

	// case 1. 创建默认的帧日志处理器
	defaultLogger := processors.NewDefaultFrameLoggerProcessor()
	// 创建另一个处理器用于链接
	nextProcessor := processors.NewLoggerProcessor("NextProcessor")
	defaultLogger.Link(nextProcessor)
	// 处理帧
	defaultLogger.ProcessFrame(frame, processors.FrameDirectionDownstream)

	// case 2. 创建自定义的帧日志处理器
	customLogger := processors.NewFrameLoggerProcessor(
		"CustomLoggerProcessor",
		"CustomPrefix",
		//[]frames.Frame{}, // 不忽略任何帧类型
		[]frames.Frame{frame}, // 忽略frame
		[]frames.Frame{},      // 包含所有帧类型
	)
	// 为自定义日志处理器链接下一个处理器
	nextProcessor2 := processors.NewLoggerProcessor("NextProcessor2")
	customLogger.Link(nextProcessor2)
	// 处理帧
	customLogger.ProcessFrame(frame, processors.FrameDirectionUpstream)
}
