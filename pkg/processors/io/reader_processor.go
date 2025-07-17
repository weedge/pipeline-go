package io

import (
	"bufio"
	"io"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

// ReaderProcessor reads from an io.Reader, creates frames, and pushes them downstream.
type ReaderProcessor struct {
	processors.BaseProcessor
	reader io.Reader
}

// NewReaderProcessor creates a new ReaderProcessor.
func NewReaderProcessor(r io.Reader) *ReaderProcessor {
	return &ReaderProcessor{
		reader: r,
	}
}

// StartReading begins reading from the reader and pushing frames.
// This should be called to start the processor. It will run until the reader returns EOF.
func (p *ReaderProcessor) StartReading() {
	scanner := bufio.NewScanner(p.reader)
	for scanner.Scan() {
		p.PushFrame(&frames.TextFrame{Text: scanner.Text()}, processors.FrameDirectionDownstream)
	}
	// When done, push an EndFrame to signal the end of the stream.
	p.PushFrame(frames.EndFrame{}, processors.FrameDirectionDownstream)
}
