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
- **Extensible**: Easily create your own custom processors by implementing the `IFrameProcessor` interface.
- **Concurrency-Safe**: Designed with concurrency in mind, using Go channels and goroutines for asynchronous processing.
- **Async Processing**: `AsyncFrameProcessor` enables asynchronous frame handling with interruption support.
- **Metrics Collection**: Enhanced processors with built-in metrics collection for TTFB and processing time.

## Directory Structure

```
/
├── pkg/
│   ├── frames/      # Defines all data and control frame types
│   ├── notifiers/   # A simple channel-based notification system
│   ├── pipeline/    # Core logic for Pipeline, PipelineTask, and parallel variants
│   └── processors/  # All built-in IFrameProcessor implementations
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
		[]processors.IFrameProcessor{
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

### Async Processor Example

Here is an example of using the `AsyncFrameProcessor` for asynchronous frame handling:

```go
package main

import (
	"log"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func main() {
	log.Println("Starting async pipeline example...")

	// 1. Create an async processor
	asyncProc := processors.NewAsyncFrameProcessor("example")

	// 2. Link it to a logger processor
	logger := processors.NewFrameTraceLogger("async", 0)
	asyncProc.Link(logger)

	// 3. Create a simple pipeline with the async processor
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			asyncProc,
		},
		nil, nil,
	)

	// 4. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	go task.Run()

	// 5. Send frames to the pipeline
	log.Println("Queueing frames...")
	task.QueueFrame(frames.NewTextFrame("Hello, async world!"))
	task.QueueFrame(frames.NewTextFrame("Processing asynchronously..."))

	// 6. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewEndFrame())

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	// Cleanup
	asyncProc.Cleanup()

	log.Println("Async pipeline finished.")
}
```

### Metrics Collection Example

Here is an example of using the enhanced processors with metrics collection:

```go
package main

import (
	"log"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// CustomProcessor is a custom processor that can generate metrics
type CustomProcessor struct {
	*processors.FrameProcessor
}

// NewCustomProcessor creates a new CustomProcessor
func NewCustomProcessor(name string) *CustomProcessor {
	return &CustomProcessor{
		FrameProcessor: processors.NewFrameProcessor(name),
	}
}

// CanGenerateMetrics returns whether the processor can generate metrics
func (p *CustomProcessor) CanGenerateMetrics() bool {
	return true
}

// ProcessFrame implements the IFrameProcessor interface
func (p *CustomProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Start metrics collection
	p.StartProcessingMetrics()
	defer p.StopProcessingMetrics()

	// Call the parent implementation
	p.FrameProcessor.ProcessFrame(frame, direction)
}

func main() {
	log.Println("Starting metrics example...")

	// 1. Create a custom processor that can generate metrics
	customProc := NewCustomProcessor("custom-processor")

	// 2. Link it to a logger processor
	logger := processors.NewFrameTraceLogger("metrics", 10) // Add small delay to see the processing
	customProc.Link(logger)

	// 3. Create a simple pipeline with the custom processor
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			customProc,
		},
		nil, nil,
	)

	// 4. Create and run a pipeline task with metrics enabled
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{
		EnableMetrics: true,
	})
	go task.Run()

	// 5. Send a start frame with metrics enabled
	startFrame := frames.NewStartFrame()
	startFrame.EnableMetrics = true
	task.QueueFrame(startFrame)

	// 6. Send frames to the pipeline
	log.Println("Queueing frames...")
	task.QueueFrame(frames.NewTextFrame("Hello, metrics world!"))
	task.QueueFrame(frames.NewTextFrame("Processing with metrics..."))

	// 7. Send a stop frame to terminate the pipeline
	task.QueueFrame(frames.NewEndFrame())

	// Give some time for processing
	time.Sleep(200 * time.Millisecond)

	// Cleanup
	customProc.Cleanup()

	log.Println("Metrics example finished.")
}
```

To run this example, you can save it as `main.go` and execute `go run .`.

## Running Tests

To run the complete test suite and see verbose output:
```bash
go test -v ./pkg/...
```

# Reference
-  pipeline-py: https://github.com/weedge/pipeline-py