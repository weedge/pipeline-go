package serializers

import (
	"encoding/json"
	"log"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

// JSONSerializer converts frames to JSON bytes.
// It ignores control frames and non-data frames.
type JSONSerializer struct {
	processors.BaseProcessor
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (s *JSONSerializer) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Ignore control frames and pass them through.
	switch frame.(type) {
	case frames.StartFrame, frames.EndFrame, frames.CancelFrame, frames.ErrorFrame, frames.StopTaskFrame, frames.MetricsFrame, frames.SyncFrame:
		s.PushFrame(frame, direction)
		return
	}

	// Attempt to marshal other frames to JSON.
	bytes, err := json.Marshal(frame)
	if err != nil {
		log.Printf("JSONSerializer: failed to marshal frame: %v", err)
		// Push an error frame instead.
		s.PushFrame(&frames.ErrorFrame{Error: err}, direction)
		return
	}

	s.PushFrame(&frames.BytesFrame{Bytes: bytes}, direction)
}

// JSONDeserializer converts JSON bytes back into a specific frame type.
type JSONDeserializer struct {
	processors.BaseProcessor
	// A factory function to create a new instance of the target frame type.
	frameFactory func() frames.Frame
}

// NewJSONDeserializer creates a new JSONDeserializer.
// The factory function should return a pointer to a new struct of the target type,
// e.g., func() frames.Frame { return &frames.TextFrame{} }.
func NewJSONDeserializer(factory func() frames.Frame) *JSONDeserializer {
	return &JSONDeserializer{
		frameFactory: factory,
	}
}

func (d *JSONDeserializer) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	bytesFrame, ok := frame.(*frames.BytesFrame)
	if !ok {
		// Not a bytes frame, so just pass it through.
		d.PushFrame(frame, direction)
		return
	}

	// Create a new frame instance to unmarshal into.
	targetFrame := d.frameFactory()
	if err := json.Unmarshal(bytesFrame.Bytes, targetFrame); err != nil {
		log.Printf("JSONDeserializer: failed to unmarshal frame: %v", err)
		d.PushFrame(&frames.ErrorFrame{Error: err}, direction)
		return
	}

	d.PushFrame(targetFrame, direction)
}
