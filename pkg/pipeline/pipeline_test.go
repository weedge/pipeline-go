package pipeline

import (
	"fmt"
	"testing"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func TestSimple(t *testing.T) {
	pipeline := NewPipeline([]processors.FrameProcessor{
		processors.NewFrameTraceLogger(),
	}, nil, nil)
	fmt.Println(pipeline)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("你好"))
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{Width: 0, Height: 0}, "PNG", "RGB"))
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}
