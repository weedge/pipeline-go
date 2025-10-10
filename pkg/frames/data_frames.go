package frames

import (
	"fmt"
	"reflect"
)

// DataFrame is a container for data frames.
type DataFrame struct {
	*BaseFrame
}

func NewDataFrame() *DataFrame {
	return &DataFrame{
		BaseFrame: NewBaseFrameWithName(reflect.TypeOf(DataFrame{}).Name()),
	}
}
func NewDataFrameWithName(name string) *DataFrame {
	return &DataFrame{
		BaseFrame: NewBaseFrameWithName(name),
	}
}

// TextFrame contains a chunk of text.
type TextFrame struct {
	*DataFrame
	Text string
}

// NewTextFrame creates a new TextFrame.
func NewTextFrame(text string) *TextFrame {
	return &TextFrame{
		DataFrame: NewDataFrameWithName("TextFrame"),
		Text:      text,
	}
}

// String returns a string representation of the TextFrame.
func (f *TextFrame) String() string {
	return fmt.Sprintf("%s(text: %s)", f.Name(), f.Text)
}

// AudioRawFrame contains a chunk of raw audio data.
type AudioRawFrame struct {
	*DataFrame
	Audio       []byte
	SampleRate  int
	NumChannels int
	SampleWidth int
	NumFrames   int
}

// NewAudioRawFrame creates a new AudioRawFrame.
func NewAudioRawFrame(audio []byte, sampleRate, numChannels, sampleWidth int) *AudioRawFrame {
	numFrames := 0
	if numChannels > 0 && sampleWidth > 0 {
		numFrames = len(audio) / (numChannels * sampleWidth)
	}
	return &AudioRawFrame{
		DataFrame:   NewDataFrameWithName("AudioRawFrame"),
		Audio:       audio,
		SampleRate:  sampleRate,
		NumChannels: numChannels,
		SampleWidth: sampleWidth,
		NumFrames:   numFrames,
	}
}

// String returns a string representation of the AudioRawFrame.
func (f *AudioRawFrame) String() string {
	return fmt.Sprintf(
		"%s(size: %d, frames: %d, sample_rate: %d, sample_width: %d, channels: %d)",
		f.Name(), len(f.Audio), f.NumFrames, f.SampleRate, f.SampleWidth, f.NumChannels,
	)
}

// ImageSize represents the dimensions of an image.
type ImageSize struct {
	Width  int
	Height int
}

// ImageRawFrame contains a raw image.
type ImageRawFrame struct {
	*DataFrame
	Image  []byte
	Size   ImageSize
	Format string
	Mode   string
}

// NewImageRawFrame creates a new ImageRawFrame.
func NewImageRawFrame(image []byte, size ImageSize, format, mode string) *ImageRawFrame {
	return &ImageRawFrame{
		DataFrame: NewDataFrameWithName("ImageRawFrame"),
		Image:     image,
		Size:      size,
		Format:    format,
		Mode:      mode,
	}
}

// String returns a string representation of the ImageRawFrame.
func (f *ImageRawFrame) String() string {
	return fmt.Sprintf(
		"%s(size: [%d, %d], format: %s, mode: %s)",
		f.Name(), f.Size.Width, f.Size.Height, f.Format, f.Mode,
	)
}

// ExtractDataFrame 从帧中提取DataFrame（通过反射检查是否嵌入了DataFrame）
func ExtractDataFrame(frame Frame) *DataFrame {
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

		if fieldType.Type == reflect.TypeOf(&DataFrame{}) {
			if !field.IsNil() {
				return field.Interface().(*DataFrame)
			}
		}
	}

	return nil
}
