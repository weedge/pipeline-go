package serializers

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/idl"
)

// ProtobufSerializer implements the Serializer interface for Protobuf.
type ProtobufSerializer struct{}

// NewProtobufSerializer creates a new ProtobufSerializer.
func NewProtobufSerializer() *ProtobufSerializer {
	return &ProtobufSerializer{}
}

// Serialize converts a frame object into a Protobuf byte slice.
func (s *ProtobufSerializer) Serialize(frame frames.Frame) ([]byte, error) {
	var pbFrame idl.Frame
	switch f := frame.(type) {
	case *frames.TextFrame:
		pbFrame.Frame = &idl.Frame_Text{
			Text: &idl.TextFrame{
				Id:   f.ID(),
				Name: f.Name(),
				Text: f.Text,
			},
		}
	case *frames.AudioRawFrame:
		pbFrame.Frame = &idl.Frame_Audio{
			Audio: &idl.AudioRawFrame{
				Id:          f.ID(),
				Name:        f.Name(),
				Audio:       f.Audio,
				SampleRate:  uint32(f.SampleRate),
				NumChannels: uint32(f.NumChannels),
				SampleWidth: uint32(f.SampleWidth),
			},
		}
	case *frames.ImageRawFrame:
		pbFrame.Frame = &idl.Frame_Image{
			Image: &idl.ImageRawFrame{
				Id:     f.ID(),
				Name:   f.Name(),
				Image:  f.Image,
				Size:   fmt.Sprintf("%dx%d", f.Size.Width, f.Size.Height),
				Format: f.Format,
				Mode:   f.Mode,
			},
		}
	default:
		return nil, fmt.Errorf("unsupported frame type for protobuf serialization: %T", f)
	}

	return proto.Marshal(&pbFrame)
}

// Deserialize converts a Protobuf byte slice back into a frame object.
func (s *ProtobufSerializer) Deserialize(data []byte) (frames.Frame, error) {
	var pbFrame idl.Frame
	if err := proto.Unmarshal(data, &pbFrame); err != nil {
		return nil, err
	}

	switch f := pbFrame.Frame.(type) {
	case *idl.Frame_Text:
		textFrame := f.Text
		return frames.NewTextFrame(textFrame.Text), nil
	case *idl.Frame_Audio:
		audioFrame := f.Audio
		return frames.NewAudioRawFrame(
			audioFrame.Audio,
			int(audioFrame.SampleRate),
			int(audioFrame.NumChannels),
			int(audioFrame.SampleWidth),
		), nil
	case *idl.Frame_Image:
		imageFrame := f.Image
		sizeParts := strings.Split(imageFrame.Size, "x")
		width, _ := strconv.Atoi(sizeParts[0])
		height, _ := strconv.Atoi(sizeParts[1])
		size := frames.ImageSize{Width: width, Height: height}
		return frames.NewImageRawFrame(
			imageFrame.Image,
			size,
			imageFrame.Format,
			imageFrame.Mode,
		), nil
	default:
		return nil, fmt.Errorf("unknown frame type in protobuf")
	}
}
