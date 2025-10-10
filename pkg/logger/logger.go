package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerConfig 日志配置结构体
type LoggerConfig struct {
	Level      slog.Level // 日志级别
	Format     string     // 日志格式 (json/text)
	Output     string     // 输出方式 (console/file/both)
	FilePath   string     // 日志文件路径
	MaxSize    int        // 单个日志文件最大大小(MB)
	MaxBackups int        // 保留旧日志文件最大个数
	MaxAge     int        // 保留旧日志文件最大天数
	Compress   bool       // 是否压缩旧日志文件
	AddSource  bool       // 是否添加源码信息
	CallerSkip int        // 调用栈跳过层数
}

func NewDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      slog.LevelInfo,
		Format:     "json",
		Output:     "console",
		FilePath:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
		AddSource:  true,
		CallerSkip: 3,
	}
}

// CallerSkipHandler wraps an existing slog.Handler and adjusts the caller depth.
type CallerSkipHandler struct {
	slog.Handler
	skip int
}

// NewCallerSkipHandler creates a new CallerSkipHandler.
func NewCallerSkipHandler(h slog.Handler, skip int) *CallerSkipHandler {
	return &CallerSkipHandler{Handler: h, skip: skip}
}

// Handle implements the slog.Handler interface.
func (h *CallerSkipHandler) Handle(ctx context.Context, r slog.Record) error {
	// Get the current call stack.
	pcs := make([]uintptr, 10)          // Adjust size as needed
	n := runtime.Callers(h.skip+2, pcs) // +2 to account for Handle and runtime.Callers itself

	// Update the record's PC to reflect the desired caller.
	if n > 0 {
		r.PC = pcs[0]
	}

	return h.Handler.Handle(ctx, r)
}

var Logger *slog.Logger = slog.New(
	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, AddSource: true,
	}),
)

// InitLoggerWithConfig 使用配置结构体初始化日志
func InitLoggerWithConfig(config *LoggerConfig) {
	var writers []io.Writer
	if config.Output == "console" || config.Output == "both" {
		writers = append(writers, os.Stdout)
	}
	if config.Output == "file" || config.Output == "both" {
		writers = append(writers, &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		})
	}
	mw := io.MultiWriter(writers...)

	// 创建 HandlerOptions
	handlerOptions := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(mw, handlerOptions)
	} else {
		handler = slog.NewTextHandler(mw, handlerOptions)
	}

	if config.CallerSkip > 0 { // 添加 caller 跳过, need AddSource: true
		handler = NewCallerSkipHandler(handler, config.CallerSkip)
	}

	Logger = slog.New(handler)
}

// InitLogger 使用独立参数初始化日志（为了向后兼容）
func InitLogger(level slog.Level, format, output, filePath string, maxSize, maxBackups, maxAge int, compress, addSource bool) {
	config := &LoggerConfig{
		Level:      level,
		Format:     format,
		Output:     output,
		FilePath:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		AddSource:  addSource,
		CallerSkip: 0,
	}
	InitLoggerWithConfig(config)
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Infof(format string, args ...any) {
	Logger.Info(fmt.Sprintf(format, args...))
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func Errorf(format string, args ...any) {
	Logger.Error(fmt.Sprintf(format, args...))
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Warnf(format string, args ...any) {
	Logger.Warn(fmt.Sprintf(format, args...))
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Debugf(format string, args ...any) {
	Logger.Debug(fmt.Sprintf(format, args...))
}
