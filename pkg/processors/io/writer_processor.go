package io

import (
	"fmt"
	"io"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

// WriterProcessor writes frames to an io.Writer.
type WriterProcessor struct {
	processors.BaseProcessor
	writer io.Writer
}

// NewWriterProcessor creates a new WriterProcessor.
func NewWriterProcessor(w io.Writer) *WriterProcessor {
	return &WriterProcessor{
		writer: w,
	}
}

// ProcessFrame writes the frame's content to the writer.
func (p *WriterProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// We only care about downstream frames for writing.
	if direction == processors.FrameDirectionUpstream {
		p.PushFrame(frame, direction)
		return
	}

	var output string
	switch f := frame.(type) {
	case *frames.TextFrame:
		output = f.Text
	case *frames.BytesFrame:
		output = string(f.Bytes)
	// Add other frame types here as needed.
	default:
		// For other types, maybe just pass them through or log them.
		p.PushFrame(frame, direction)
		return
	}

	// Write the output followed by a newline.
	_, err := fmt.Fprintln(p.writer, output)
	if err != nil {
		// Handle the error, maybe push an ErrorFrame.
		errFrame := &frames.ErrorFrame{Error: err, Fatal: false}
		p.PushFrame(errFrame, processors.FrameDirectionUpstream) // Push error upstream
	}
}
