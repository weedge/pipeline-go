package main

import (
	"log"
	"time"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
	"github.com/wuyong/pipeline-go/pkg/serializers"
)

func main() {
	log.Println("Starting serializer example")

	// 1. Create the serializer, deserializer, and a logger
	serializer := serializers.NewJSONSerializer()
	deserializer := serializers.NewJSONDeserializer(func() frames.Frame { return &frames.TextFrame{} })
	logger := processors.NewLoggerProcessor("OutputLogger")

	// 2. Create a pipeline for a serialization round trip
	pl := pipeline.NewPipeline(
		[]processors.FrameProcessor{serializer, deserializer, logger},
		nil, nil,
	)

	// 3. Create a pipeline task
	params := pipeline.PipelineParams{}
	task := pipeline.NewPipelineTask(pl, params)

	// 4. Run the task in a separate goroutine
	go task.Run()

	// 5. Queue a text frame
	log.Println("Queueing frame")
	task.QueueFrame(&frames.TextFrame{Text: "serialize me"})
	task.QueueFrame(frames.EndFrame{})

	// Wait for the task to finish
	for !task.HasFinished() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Serializer example finished")
}
