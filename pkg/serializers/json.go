package serializers

import (
	"encoding/json"
	"fmt"

	"github.com/wuyong/pipeline-go/pkg/frames"
)

const (
	frameTypeText  = "text"
	frameTypeAudio = "audio"
	frameTypeImage = "image"
)

// jsonFrameWrapper is used to wrap frame data with type information for JSON serialization.
type jsonFrameWrapper struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// JsonSerializer implements the Serializer interface for JSON.
type JsonSerializer struct{}

// NewJsonSerializer creates a new JsonSerializer.
func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

// Serialize converts a frame object into a JSON byte slice.
func (s *JsonSerializer) Serialize(frame frames.Frame) ([]byte, error) {
	var frameType string
	switch frame.(type) {
	case *frames.TextFrame:
		frameType = frameTypeText
	case *frames.AudioRawFrame:
		frameType = frameTypeAudio
	case *frames.ImageRawFrame:
		frameType = frameTypeImage
	default:
		return nil, fmt.Errorf("unsupported frame type for json serialization: %T", frame)
	}

	// A temporary struct to hold the type and data for serialization.
	typedFrame := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: frameType,
		Data: frame,
	}

	return json.Marshal(typedFrame)
}

// Deserialize converts a JSON byte slice back into a frame object.
func (s *JsonSerializer) Deserialize(data []byte) (frames.Frame, error) {
	var wrapper jsonFrameWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("error unmarshalling frame wrapper: %w", err)
	}

	switch wrapper.Type {
	case frameTypeText:
		var frame frames.TextFrame
		if err := json.Unmarshal(wrapper.Data, &frame); err != nil {
			return nil, fmt.Errorf("error unmarshalling text frame: %w", err)
		}
		return &frame, nil
	case frameTypeAudio:
		var frame frames.AudioRawFrame
		if err := json.Unmarshal(wrapper.Data, &frame); err != nil {
			return nil, fmt.Errorf("error unmarshalling audio frame: %w", err)
		}
		return &frame, nil
	case frameTypeImage:
		var frame frames.ImageRawFrame
		if err := json.Unmarshal(wrapper.Data, &frame); err != nil {
			return nil, fmt.Errorf("error unmarshalling image frame: %w", err)
		}
		return &frame, nil
	default:
		return nil, fmt.Errorf("unknown frame type in json: %s", wrapper.Type)
	}
}
