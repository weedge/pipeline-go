package frames

import (
	"fmt"
	"reflect"
)

// SystemFrame is a frame for system-level events.
type SystemFrame struct {
	*BaseFrame
}

func NewSystemFrame() *SystemFrame {
	return &SystemFrame{
		BaseFrame: NewBaseFrame(reflect.TypeOf(SystemFrame{}).Name()),
	}
}

// CancelFrame indicates that a pipeline needs to stop right away.
type CancelFrame struct {
	*SystemFrame
}

func NewCancelFrame() *CancelFrame {
	return &CancelFrame{
		SystemFrame: NewSystemFrame(),
	}
}

// ErrorFrame is used to notify upstream that an error has occurred.
type ErrorFrame struct {
	*SystemFrame
	Error error
	Fatal bool
}

func NewErrorFrame(err error, fatal bool) *ErrorFrame {
	return &ErrorFrame{
		SystemFrame: NewSystemFrame(),
		Error:       err,
		Fatal:       fatal,
	}
}

func (f *ErrorFrame) String() string {
	return fmt.Sprintf("%s(error: %s, fatal: %t)", f.Name(), f.Error, f.Fatal)
}

// StopTaskFrame indicates that a pipeline task should be stopped.
type StopTaskFrame struct {
	*SystemFrame
}

func NewStopTaskFrame() *StopTaskFrame {
	return &StopTaskFrame{
		SystemFrame: NewSystemFrame(),
	}
}

// StartInterruptionFrame indicates that a user has started speaking.
type StartInterruptionFrame struct {
	*SystemFrame
}

func NewStartInterruptionFrame() *StartInterruptionFrame {
	return &StartInterruptionFrame{
		SystemFrame: NewSystemFrame(),
	}
}

// StopInterruptionFrame indicates that a user has stopped speaking.
type StopInterruptionFrame struct {
	*SystemFrame
}

func NewStopInterruptionFrame() *StopInterruptionFrame {
	return &StopInterruptionFrame{
		SystemFrame: NewSystemFrame(),
	}
}

// MetricsFrame contains metrics about the pipeline.
type MetricsFrame struct {
	*SystemFrame
	TTFB       []map[string]interface{}
	Processing []map[string]interface{}
	Tokens     []map[string]interface{}
	Characters []map[string]interface{}
}

func NewMetricsFrame() *MetricsFrame {
	return &MetricsFrame{
		SystemFrame: NewSystemFrame(),
	}
}

func (f *MetricsFrame) String() string {
	return fmt.Sprintf("%s ttfb:%v | processing:%v | tokens:%v | characters:%v",
		f.Name(), f.TTFB, f.Processing, f.Tokens, f.Characters)
}

// UsageMetricFrame contains usage metrics.
type UsageMetricFrame struct {
	*SystemFrame
	Key   string
	Value int
}

func NewUsageMetricFrame(key string, value int) *UsageMetricFrame {
	return &UsageMetricFrame{
		SystemFrame: NewSystemFrame(),
		Key:         key,
		Value:       value,
	}
}

func (f *UsageMetricFrame) String() string {
	return fmt.Sprintf("%s(key: %s, value: %d)", f.Name(), f.Key, f.Value)
}
