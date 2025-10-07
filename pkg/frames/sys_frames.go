package frames

import (
	"fmt"
)

// SystemFrame is a frame for system-level events.
type SystemFrame struct {
	*BaseFrame
}

func NewSystemFrame() *SystemFrame {
	return &SystemFrame{
		BaseFrame: NewBaseFrameWithName("SystemFrame"),
	}
}

// CancelFrame indicates that a pipeline needs to stop right away.
type CancelFrame struct {
	*SystemFrame
}

func NewCancelFrame() *CancelFrame {
	return &CancelFrame{
		SystemFrame: &SystemFrame{
			BaseFrame: NewBaseFrameWithName("CancelFrame"),
		},
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
		SystemFrame: &SystemFrame{
			BaseFrame: NewBaseFrameWithName("ErrorFrame"),
		},
		Error: err,
		Fatal: fatal,
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
		SystemFrame: &SystemFrame{
			BaseFrame: NewBaseFrameWithName("StopTaskFrame"),
		},
	}
}

// StartInterruptionFrame indicates that a user has started speaking.
type StartInterruptionFrame struct {
	*SystemFrame
}

func NewStartInterruptionFrame() *StartInterruptionFrame {
	return &StartInterruptionFrame{
		SystemFrame: &SystemFrame{
			BaseFrame: NewBaseFrameWithName("StartInterruptionFrame"),
		},
	}
}

// StopInterruptionFrame indicates that a user has stopped speaking.
type StopInterruptionFrame struct {
	*SystemFrame
}

func NewStopInterruptionFrame() *StopInterruptionFrame {
	return &StopInterruptionFrame{
		SystemFrame: &SystemFrame{
			BaseFrame: NewBaseFrameWithName("StopInterruptionFrame"),
		},
	}
}

// MetricsFrame contains metrics about the pipeline.
type MetricsFrame struct {
	*SystemFrame
	TTFB       []map[string]any
	Processing []map[string]any
	Tokens     []map[string]any
	Characters []map[string]any
}

func NewMetricsFrame() *MetricsFrame {
	return &MetricsFrame{
		SystemFrame: NewSystemFrame(),
	}
}

// NewMetricsFrameWithTTFB creates a new MetricsFrame with TTFB metrics.
func NewMetricsFrameWithTTFB(ttfb []map[string]interface{}) *MetricsFrame {
	return &MetricsFrame{
		SystemFrame: &SystemFrame{
			BaseFrame: NewBaseFrameWithName("MetricsFrame"),
		},
		TTFB: ttfb,
	}
}

// NewMetricsFrameWithProcessing creates a new MetricsFrame with processing metrics.
func NewMetricsFrameWithProcessing(processing []map[string]interface{}) *MetricsFrame {
	return &MetricsFrame{
		SystemFrame: &SystemFrame{
			BaseFrame: NewBaseFrameWithName("MetricsFrame"),
		},
		Processing: processing,
	}
}

func (f *MetricsFrame) String() string {
	return fmt.Sprintf("%s ttfb:%+v | processing:%+v | tokens:%+v | characters:%+v",
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
