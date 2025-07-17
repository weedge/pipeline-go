package pipeline

import (
	"fmt"
	"testing"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func TestSimple(t *testing.T) {
	pipeline := NewPipeline([]processors.FrameProcessor{
		processors.NewFrameTraceLogger("simple", 0),
	}, nil, nil)
	fmt.Println(pipeline)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("你好"))
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{Width: 0, Height: 0}, "PNG", "RGB"))
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}

func TestParallelPipeline(t *testing.T) {
	pipeline := NewPipeline([]processors.FrameProcessor{
		processors.NewFrameTraceLogger("0", 0),
		NewParallelPipeline(
			[]processors.FrameProcessor{processors.NewFrameTraceLogger("1.0", 1000)},
			[]processors.FrameProcessor{processors.NewFrameTraceLogger("1.1", 0)},
		),
		processors.NewFrameTraceLogger("3", 0),
	}, nil, nil)
	fmt.Println(pipeline)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("你好"))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}

func TestSyncParallelPipeline(t *testing.T) {
	pipeline := NewPipeline([]processors.FrameProcessor{
		processors.NewFrameTraceLogger("0", 0),
		NewSyncParallelPipeline(
			NewPipeline([]processors.FrameProcessor{processors.NewFrameTraceLogger("1.0", 1000)}, nil, nil),
			NewPipeline([]processors.FrameProcessor{processors.NewFrameTraceLogger("1.1", 0)}, nil, nil),
		),
		processors.NewFrameTraceLogger("3", 0),
	}, nil, nil)
	fmt.Println(pipeline)

	task := NewPipelineTask(pipeline, PipelineParams{})

	task.QueueFrame(frames.NewTextFrame("你好"))
	task.QueueFrame(frames.NewEndFrame())

	task.Run()
}
