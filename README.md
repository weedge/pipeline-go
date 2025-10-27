# Pipeline-Go

Pipeline-Go is a Go library for building flexible, concurrent data processing pipelines. It is inspired by a similar Python project and provides a structured way to process streams of data "frames" through a series of composable "processors".

## Design
⭐️ [pipeline design](https://github.com/ai-bot-pro/pipeline-py#design) ⭐️

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
    - **Output Processing**: Simple `OutputProcessor` for basic frame handling, and advanced `AdvancedOutputProcessor` and `OutputFrameProcessor` for complex output scenarios with async processing, interruption handling, and metrics support.
- **Extensible**: Easily create your own custom processors by implementing the `IFrameProcessor` interface.
- **Concurrency-Safe**: Designed with concurrency in mind, using Go channels and goroutines for asynchronous processing.
- **Async Processing**: `AsyncFrameProcessor` enables asynchronous frame handling with interruption support.
- **Metrics Collection**: Enhanced processors with built-in metrics collection for TTFB and processing time.

## Directory Structure

```
/
├── examples/        # processor examples
├── pkg/
│   ├── logger/      # defined logger to config
│   ├── frames/      # Defines all data and control frame types
│   ├── notifiers/   # A simple channel-based notification system
│   ├── pipeline/    # Core logic for Pipeline, PipelineTask, and parallel variants
│   ├── idl/         # rpc IDL
│   ├── serializers/ # pb json serializers
│   └── processors/  # All built-in IFrameProcessor implementations
├── go.mod
```

## Usage Example
see [examples](./examples/)

### Filter Processor Example
```shell
go run examples/filter_example.go
```

### Async Processor Example

```shell
go run examples/async_example.go
# with interrupt
go run examples/async_interruption_example.go
```

### Advanced Output Processor Example

```shell
go run examples/advanced_output_example.go
```

### Metrics Collection Example

```shell
go run examples/metrics_example.go
```

## Running Tests

To run the complete test suite and see verbose output:
```bash
go test -v ./pkg/...
```

# Reference
-  pipeline-py: https://github.com/weedge/pipeline-py
