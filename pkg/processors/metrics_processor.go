package processors

import (
	"log"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
)

// MetricsProcessor handles processor performance metrics (TTFB and processing time).
type MetricsProcessor struct {
	name                  string
	startTTFBTime         time.Time
	startProcessingTime   time.Time
	shouldReportTTFB      bool
	reportOnlyInitialTTFB bool
}

// NewMetricsProcessor creates a new MetricsProcessor.
func NewMetricsProcessor(name string) *MetricsProcessor {
	return &MetricsProcessor{
		name:             name,
		shouldReportTTFB: true,
	}
}

// StartTTFBMetrics starts TTFB metrics collection.
func (m *MetricsProcessor) StartTTFBMetrics(reportOnlyInitialTTFB bool) {
	if m.shouldReportTTFB {
		m.startTTFBTime = time.Now()
		m.reportOnlyInitialTTFB = reportOnlyInitialTTFB
		m.shouldReportTTFB = !reportOnlyInitialTTFB
	}
}

// StopTTFBMetrics stops TTFB metrics collection and returns a MetricsFrame.
func (m *MetricsProcessor) StopTTFBMetrics() *frames.MetricsFrame {
	if m.startTTFBTime.IsZero() {
		return nil
	}

	value := time.Since(m.startTTFBTime).Seconds()
	log.Printf("%s TTFB: %f", m.name, value)
	ttfb := map[string]interface{}{"processor": m.name, "value": value}
	m.startTTFBTime = time.Time{} // Reset to zero time
	return frames.NewMetricsFrameWithTTFB([]map[string]interface{}{ttfb})
}

// StartProcessingMetrics starts processing time metrics collection.
func (m *MetricsProcessor) StartProcessingMetrics() {
	m.startProcessingTime = time.Now()
}

// StopProcessingMetrics stops processing time metrics collection and returns a MetricsFrame.
func (m *MetricsProcessor) StopProcessingMetrics() *frames.MetricsFrame {
	if m.startProcessingTime.IsZero() {
		return nil
	}

	value := time.Since(m.startProcessingTime).Seconds()
	log.Printf("%s processing time: %f", m.name, value)
	processing := map[string]interface{}{"processor": m.name, "value": value}
	m.startProcessingTime = time.Time{} // Reset to zero time
	return frames.NewMetricsFrameWithProcessing([]map[string]interface{}{processing})
}
