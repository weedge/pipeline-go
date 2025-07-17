# Go Pipeline

This is a Go implementation of a generic, asynchronous pipeline for processing frames of data. It is a refactoring of the original Python project.

## Features

- Generic `FrameProcessor` interface
- Asynchronous, channel-based `Pipeline` and `PipelineTask`
- Protobuf-defined data frames for cross-compatibility
- Simple, composable design

## Installation

To use this library in your project:
```bash
go get github.com/weedge/pipeline-go
```

## Usage

Below is a simple example of creating and running a pipeline. For a more detailed example, see `examples/sentence_aggregator`.

```go
package main

import (
	"context"
	"fmt"
	"github.com/weedge/pipeline-go/frames"
	"github.com/weedge/pipeline-go/pipeline"
	"github.com/weedge/pipeline-go/processors"
)

// SimpleProcessor is a basic example of a frame processor.
type SimpleProcessor struct {
	pipeline.BaseProcessor
}

func (p *SimpleProcessor) ProcessFrame(frame frames.PipelineFrame, direction processors.FrameDirection) error {
	fmt.Printf("Processing frame: %s\n", frame.GetFrameName())
	return p.BaseProcessor.ProcessFrame(frame, direction)
}

func main() {
	// 1. Create a processor
	simpleProcessor := &SimpleProcessor{}

	// 2. Create a new pipeline
	myPipeline := pipeline.NewPipeline(simpleProcessor)

	// 3. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline)
	go task.Run(context.Background())

	// 4. Send a frame
	task.InChan() <- &frames.TextFrame{Text: "Hello, world!"}

	// ... wait for processing or handle output ...
}
```

## Running Tests

To run the test suite:
```bash
go test ./... -v
```

## Running Examples

To run the sentence aggregator example:
```bash
go run ./examples/sentence_aggregator/main.go
```
```