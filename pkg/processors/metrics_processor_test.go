package processors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weedge/pipeline-go/pkg/frames"
)

func TestMetricsProcessor(t *testing.T) {
	// Create a new metrics processor
	metrics := NewMetricsProcessor("test-processor")

	// Test TTFB metrics
	metrics.StartTTFBMetrics(false)
	time.Sleep(10 * time.Millisecond)
	frame := metrics.StopTTFBMetrics()

	assert.NotNil(t, frame)
	assert.NotEmpty(t, frame.TTFB)
	assert.Equal(t, "test-processor", frame.TTFB[0]["processor"])
	assert.Greater(t, frame.TTFB[0]["value"], float64(0))

	// Test that stopping without starting returns nil
	frame = metrics.StopTTFBMetrics()
	assert.Nil(t, frame)

	// Test processing metrics
	metrics.StartProcessingMetrics()
	time.Sleep(10 * time.Millisecond)
	frame = metrics.StopProcessingMetrics()

	assert.NotNil(t, frame)
	assert.NotEmpty(t, frame.Processing)
	assert.Equal(t, "test-processor", frame.Processing[0]["processor"])
	assert.Greater(t, frame.Processing[0]["value"], float64(0))

	// Test that stopping without starting returns nil
	frame = metrics.StopProcessingMetrics()
	assert.Nil(t, frame)
}

func TestFrameProcessor(t *testing.T) {
	// Create a new enhanced processor
	processor := NewFrameProcessor("test-processor")

	// Test basic properties
	assert.Equal(t, "test-processor", processor.Name())
	assert.False(t, processor.InterruptionsAllowed())
	assert.False(t, processor.MetricsEnabled())
	assert.False(t, processor.UsageMetricsEnabled())
	assert.False(t, processor.ReportOnlyInitialTTFB())

	// Test that by default it cannot generate metrics
	assert.False(t, processor.CanGenerateMetrics())

	// Test skip frames functionality
	frame := frames.NewTextFrame("test")
	processor.AddSkipFrame(frame)

	// We can't easily test the skip functionality without a full pipeline,
	// but we can verify the frame was added to the skip list
	// This would normally be tested in an integration test

	// Test linking
	nextProcessor := NewFrameProcessor("next-processor")
	processor.Link(nextProcessor)

	// Test string representation
	assert.Equal(t, "test-processor", processor.String())
}

func TestFrameProcessorStartFrame(t *testing.T) {
	// Create a new enhanced processor
	processor := NewFrameProcessor("test-processor")

	// Create a start frame with specific settings using the constructor
	startFrame := frames.NewStartFrame()
	startFrame.AllowInterruptions = true
	startFrame.EnableMetrics = true
	startFrame.EnableUsageMetrics = true
	startFrame.ReportOnlyInitialTTFB = true

	// Process the start frame
	processor.ProcessFrame(startFrame, FrameDirectionDownstream)

	// Check that properties were set correctly
	assert.True(t, processor.allowInterruptions)
	assert.True(t, processor.enableMetrics)
	assert.True(t, processor.enableUsageMetrics)
	assert.True(t, processor.reportOnlyInitialTTFB)
}

func TestFrameProcessorMetrics(t *testing.T) {
	// Create a new frame processor
	processor := NewFrameProcessor("test-processor")

	// Mock the CanGenerateMetrics method to return true
	// Since we can't easily override methods in Go, we'll test the metrics functions directly
	// by calling them even though CanGenerateMetrics returns false

	// Enable metrics through a start frame
	startFrame := &frames.StartFrame{
		EnableMetrics: true,
	}
	processor.ProcessFrame(startFrame, FrameDirectionDownstream)

	// Test TTFB metrics
	processor.StartTTFBMetrics()
	time.Sleep(10 * time.Millisecond)
	processor.StopTTFBMetrics()

	// Test processing metrics
	processor.StartProcessingMetrics()
	time.Sleep(10 * time.Millisecond)
	processor.StopProcessingMetrics()

	// Test stopping all metrics
	processor.StopAllMetrics()
}
