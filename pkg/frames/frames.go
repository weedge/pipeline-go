package frames

// Frame is the interface for all frames that flow through the pipeline.
type Frame interface{}

// StartFrame is a control frame that signals the start of the pipeline.
type StartFrame struct {
	AllowInterruptions     bool
	EnableMetrics          bool
	EnableUsageMetrics     bool
	ReportOnlyInitialTTFB  bool
}

// EndFrame is a control frame that signals the end of the pipeline.
type EndFrame struct{}

// CancelFrame is a control frame that signals the cancellation of the pipeline.
type CancelFrame struct{}

// ErrorFrame is a control frame that signals an error in the pipeline.
type ErrorFrame struct {
	Error error
	Fatal bool
}

// StopTaskFrame is a control frame that signals the task to stop.
type StopTaskFrame struct{}

// MetricsFrame is a frame that contains metrics data.
type MetricsFrame struct {
	TTFB       []map[string]interface{}
	Processing []map[string]interface{}
}

// SyncFrame is a control frame used for synchronization.
type SyncFrame struct{}

// TextFrame is a data frame containing text.
type TextFrame struct {
	Text string
}

// BytesFrame is a data frame containing raw bytes.
type BytesFrame struct {
	Bytes []byte
}
