package frames

import (
	"fmt"
	"reflect"
)

// ControlFrame is a frame for controlling the pipeline.
type ControlFrame struct {
	*BaseFrame
}

func NewControlFrame() *ControlFrame {
	return &ControlFrame{
		BaseFrame: NewBaseFrameWithName("ControlFrame"),
	}
}

func NewControlFrameWithName(name string) *ControlFrame {
	return &ControlFrame{
		BaseFrame: NewBaseFrameWithName(name),
	}
}

// StartFrame is the first frame that should be pushed down a pipeline.
type StartFrame struct {
	*ControlFrame
	AllowInterruptions    bool
	EnableMetrics         bool
	EnableUsageMetrics    bool
	ReportOnlyInitialTTFB bool
	AudioInSampleRate     int
	AudioOutSampleRate    int
	IsPushBlock           bool
	IsUpPushBlock         bool
}

// new default start frame
func NewStartFrame() *StartFrame {
	return &StartFrame{
		ControlFrame: &ControlFrame{
			BaseFrame: NewBaseFrameWithName("StartFrame"),
		},
		AllowInterruptions:    false,
		EnableMetrics:         false,
		EnableUsageMetrics:    false,
		ReportOnlyInitialTTFB: false,
		AudioInSampleRate:     16000,
		AudioOutSampleRate:    16000,
		IsPushBlock:           false,
		IsUpPushBlock:         false,
	}
}
func (f *StartFrame) String() string {
	return fmt.Sprintf("%s(allowInterruptions: %t, enableMetrics: %t, enableUsageMetrics: %t, reportOnlyInitialTTFB: %t, audioInSampleRate: %d, audioOutSampleRate: %d, IsPushBlock: %t, IsUpPushBlock: %t)", f.Name(), f.AllowInterruptions, f.EnableMetrics, f.EnableUsageMetrics, f.ReportOnlyInitialTTFB, f.AudioInSampleRate, f.AudioOutSampleRate, f.IsPushBlock, f.IsUpPushBlock)
}

// EndFrame indicates that a pipeline has ended.
type EndFrame struct {
	*ControlFrame
}

func NewEndFrame() *EndFrame {
	return &EndFrame{
		ControlFrame: &ControlFrame{
			BaseFrame: NewBaseFrameWithName("EndFrame"),
		},
	}
}

// SyncFrame is used to know when the internal pipelines have finished.
type SyncFrame struct {
	*ControlFrame
}

func NewSyncFrame() *SyncFrame {
	return &SyncFrame{
		ControlFrame: &ControlFrame{
			BaseFrame: NewBaseFrameWithName("SyncFrame"),
		},
	}
}

// SyncNotifyFrame is used to know when notification has been received.
type SyncNotifyFrame struct {
	*ControlFrame
}

func NewSyncNotifyFrame() *SyncNotifyFrame {
	return &SyncNotifyFrame{
		ControlFrame: &ControlFrame{
			BaseFrame: NewBaseFrameWithName("SyncNotifyFrame"),
		},
	}
}

// IdleFrame indicates that a pipeline has ended.
type IdleFrame struct {
	*ControlFrame
}

func NewIdleFrame() *IdleFrame {
	return &IdleFrame{
		ControlFrame: &ControlFrame{
			BaseFrame: NewBaseFrameWithName("IdleFrame"),
		},
	}
}

// extractControlFrame 从帧中提取ControlFrame（通过反射检查是否嵌入了ControlFrame）
func ExtractControlFrame(frame Frame) *ControlFrame {
	val := reflect.ValueOf(frame)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		if fieldType.Type == reflect.TypeOf(&ControlFrame{}) {
			if !field.IsNil() {
				return field.Interface().(*ControlFrame)
			}
		}
	}

	return nil
}
