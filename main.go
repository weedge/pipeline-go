package main

import (
	"bytes"
	"log"
	"strings"

	"github.com/wuyong/pipeline-go/pkg/pipeline"
	"github.com/wuyong/pipeline-go/pkg/processors"
	"github.com/wuyong/pipeline-go/pkg/processors/io"
)

func main() {
	log.Println("Starting I/O processor example")

	// 1. Prepare the input and output buffers
	inputData := "hello world\nthis is a test\n"
	inputBuffer := bytes.NewBufferString(inputData)
	var outputBuffer bytes.Buffer

	// 2. Create the I/O processors and a transformer
	reader := io.NewReaderProcessor(inputBuffer)
	transformer := processors.NewStatelessTextTransformer(strings.ToUpper)
	writer := io.NewWriterProcessor(&outputBuffer)

	// 3. Create a pipeline
	// The reader is the source, so it's not part of the pipeline itself.
	// It will push frames into the pipeline.
	pl := pipeline.NewPipeline(
		[]processors.FrameProcessor{transformer, writer},
		nil, nil,
	)
	reader.Link(pl)

	// 4. Start the reader. This will drive the pipeline.
	reader.StartReading()

	// 5. Print the result from the output buffer
	log.Printf("Pipeline finished. Output:\n%s", outputBuffer.String())
}
