package serializers

import (
	"reflect"
	"testing"

	"github.com/weedge/pipeline-go/pkg/frames"
)

func TestSerializers(t *testing.T) {
	// Sample frames to be used in tests
	testFrames := map[string]frames.Frame{
		"TextFrame":     frames.NewTextFrame("hello world"),
		"AudioRawFrame": frames.NewAudioRawFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}, 16000, 1, 2),
		"ImageRawFrame": frames.NewImageRawFrame(
			[]byte{10, 20, 30},
			frames.ImageSize{Width: 1920, Height: 1080},
			"jpeg",
			"RGB",
		),
	}

	// Serializers to be tested
	serializers := map[string]Serializer{
		"Protobuf": NewProtobufSerializer(),
		"JSON":     NewJsonSerializer(),
	}

	for serName, serializer := range serializers {
		t.Run(serName, func(t *testing.T) {
			for frameName, originalFrame := range testFrames {
				t.Run(frameName, func(t *testing.T) {
					// Serialize the frame
					serializedData, err := serializer.Serialize(originalFrame)
					if err != nil {
						t.Fatalf("Serialize() error = %v", err)
					}
					if serializedData == nil {
						t.Fatalf("Serialize() returned nil data")
					}

					// Deserialize the data
					deserializedFrame, err := serializer.Deserialize(serializedData)
					if err != nil {
						t.Fatalf("Deserialize() error = %v", err)
					}
					if deserializedFrame == nil {
						t.Fatalf("Deserialize() returned nil frame")
					}

					// Compare the meaningful fields of the frames, ignoring the base fields (ID, Name)
					// which are regenerated on deserialization.
					if !areFramesEqual(originalFrame, deserializedFrame) {
						t.Errorf("frames are not equal after serialization/deserialization cycle")
						t.Logf("Original:   %#v", originalFrame)
						t.Logf("Deserialized: %#v", deserializedFrame)
					}
				})
			}
		})
	}
}

// areFramesEqual compares the data fields of two frames, ignoring the base fields.
func areFramesEqual(a, b frames.Frame) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	switch fa := a.(type) {
	case *frames.TextFrame:
		fb := b.(*frames.TextFrame)
		return fa.Text == fb.Text
	case *frames.AudioRawFrame:
		fb := b.(*frames.AudioRawFrame)
		return reflect.DeepEqual(fa.Audio, fb.Audio) &&
			fa.SampleRate == fb.SampleRate &&
			fa.NumChannels == fb.NumChannels &&
			fa.SampleWidth == fb.SampleWidth
	case *frames.ImageRawFrame:
		fb := b.(*frames.ImageRawFrame)
		return reflect.DeepEqual(fa.Image, fb.Image) &&
			fa.Size == fb.Size &&
			fa.Format == fb.Format &&
			fa.Mode == fb.Mode
	}
	return false
}
