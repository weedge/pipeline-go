package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"

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
	}
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
	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: config.Level, AddSource: config.AddSource})
	} else {
		handler = slog.NewTextHandler(mw, &slog.HandlerOptions{Level: config.Level, AddSource: config.AddSource})
	}
	Logger = slog.New(handler)
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
