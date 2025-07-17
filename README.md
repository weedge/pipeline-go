# Pipeline-Go

Pipeline-Go is a Go library for building flexible, concurrent data processing pipelines. It is inspired by a similar Python project and provides a structured way to process streams of data "frames" through a series of composable "processors".

## Core Concepts

- **Frame**: The basic unit of data that flows through a pipeline. Frames can be of various types, such as `TextFrame`, `ImageRawFrame`, or control frames like `EndFrame` and `StartFrame`.
- **Processor**: An interface for components that act on frames. Processors are the building blocks of a pipeline. They can be used to filter, transform, aggregate, or perform any other operation on frames.
- **Pipeline**: A sequence of processors linked together. Pipelines can be simple linear sequences or more complex parallel structures.
- **Task**: An executable instance of a pipeline. A `PipelineTask` manages the lifecycle of a pipeline, including running it, queueing frames, and handling shutdown.

## Features

- **Sequential & Parallel Pipelines**: Build simple linear pipelines (`NewPipeline`) or complex concurrent pipelines that process frames in parallel (`NewParallelPipeline`, `NewSyncParallelPipeline`).
- **Rich Set of Processors**: Includes a variety of built-in processors for common tasks:
    - **Filtering**: `FrameFilter` (based on a function) and `TypeFilter` (based on frame type).
    - **Aggregation**: `SentenceAggregator`, `GatedAggregator`, `HoldFramesAggregator`, and `HoldLastFrameAggregator`.
- **Extensible**: Easily create your own custom processors by implementing the `FrameProcessor` interface.
- **Concurrency-Safe**: Designed with concurrency in mind, using Go channels and goroutines for asynchronous processing.

## Directory Structure

```
/
├── pkg/
│   ├── frames/      # Defines all data and control frame types
│   ├── notifiers/   # A simple channel-based notification system
│   ├── pipeline/    # Core logic for Pipeline, PipelineTask, and parallel variants
│   └── processors/  # All built-in FrameProcessor implementations
├── go.mod
└── main.go        # An example executable
```

## Usage Example

Here is a simple example of building a pipeline that filters out `ImageRawFrame`s and logs the remaining frames.

```go
package main

import (
	"log"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
	"github.com/weedge/pipeline-go/pkg/processors/filters"
)

func main() {
	log.Println("Starting pipeline example...")

	// 1. Create the processors
	// This filter will only allow frames that are NOT ImageRawFrames to pass
	imageFilter := filters.NewFrameFilter(func(frame frames.Frame) bool {
		if _, ok := frame.(*frames.ImageRawFrame); ok {
			return false // Drop the frame
		}
		return true // Keep the frame
	})

	// This logger will print any frame that it receives
	logger := processors.NewFrameTraceLogger("output", 0)

	// 2. Create a new pipeline
	myPipeline := pipeline.NewPipeline(
		[]processors.FrameProcessor{
			imageFilter,
			logger,
		},
		nil, nil,
	)

	// 3. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	go task.Run()

	// 4. Send frames to the pipeline
	log.Println("Queueing TextFrame...")
	task.QueueFrame(frames.NewTextFrame("Hello, world!"))

	log.Println("Queueing ImageRawFrame (will be filtered out)...")
	task.QueueFrame(frames.NewImageRawFrame([]byte{}, frames.ImageSize{}, "PNG", "RGB"))

	log.Println("Queueing AudioRawFrame...")
	task.QueueFrame(frames.NewAudioRawFrame([]byte{}, 16000, 1, 2))

	// 5. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewStopTaskFrame())

	log.Println("Pipeline finished.")
}
```

To run this example, you can save it as `main.go` and execute `go run .`.

## Running Tests

To run the complete test suite and see verbose output:
```bash
go test -v ./...
```

# Reference
-  pipeline-py: https://github.com/weedge/pipeline-py