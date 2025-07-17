package serializers

import "github.com/weedge/pipeline-go/pkg/frames"

// Serializer defines the interface for serializing and deserializing frames.
type Serializer interface {
	// Serialize converts a frame object into a byte slice.
	Serialize(frame frames.Frame) ([]byte, error)
	// Deserialize converts a byte slice back into a frame object.
	Deserialize(data []byte) (frames.Frame, error)
}
