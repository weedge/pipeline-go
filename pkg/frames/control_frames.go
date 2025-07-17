package frames

import "reflect"

// ControlFrame is a frame for controlling the pipeline.
type ControlFrame struct {
	*BaseFrame
}

func NewControlFrame() *ControlFrame {
	return &ControlFrame{
		BaseFrame: NewBaseFrame(reflect.TypeOf(ControlFrame{}).Name()),
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
}

func NewStartFrame() *StartFrame {
	return &StartFrame{
		ControlFrame:       NewControlFrame(),
		AudioInSampleRate:  16000,
		AudioOutSampleRate: 24000,
	}
}

// EndFrame indicates that a pipeline has ended.
type EndFrame struct {
	*ControlFrame
}

func NewEndFrame() *EndFrame {
	return &EndFrame{
		ControlFrame: NewControlFrame(),
	}
}

// SyncFrame is used to know when the internal pipelines have finished.
type SyncFrame struct {
	*ControlFrame
}

func NewSyncFrame() *SyncFrame {
	return &SyncFrame{
		ControlFrame: NewControlFrame(),
	}
}

// SyncNotifyFrame is used to know when notification has been received.
type SyncNotifyFrame struct {
	*ControlFrame
}

func NewSyncNotifyFrame() *SyncNotifyFrame {
	return &SyncNotifyFrame{
		ControlFrame: NewControlFrame(),
	}
}

// IdleFrame indicates that a pipeline has ended.
type IdleFrame struct {
	*ControlFrame
}

func NewIdleFrame() *IdleFrame {
	return &IdleFrame{
		ControlFrame: NewControlFrame(),
	}
}
