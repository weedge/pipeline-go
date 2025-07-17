package pipeline

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/wuyong/pipeline-go/pkg/frames"
	"github.com/wuyong/pipeline-go/pkg/processors"
)

// PipelineParams holds parameters for a pipeline task.
type PipelineParams struct {
	AllowInterruptions      bool
	EnableMetrics           bool
	EnableUsageMetrics      bool
	SendInitialEmptyMetrics bool
	ReportOnlyInitialTTFB   bool
}

// Source is a processor that handles upstream frames for a task.
type TaskSource struct {
	processors.BaseProcessor
	upQueue chan frames.Frame
}

func NewTaskSource(upQueue chan frames.Frame) *TaskSource {
	return &TaskSource{
		upQueue: upQueue,
	}
}

func (s *TaskSource) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	switch direction {
	case processors.FrameDirectionUpstream:
		s.handleUpstreamFrame(frame)
	case processors.FrameDirectionDownstream:
		s.PushFrame(frame, direction)
	}
}

func (s *TaskSource) handleUpstreamFrame(frame frames.Frame) {
	if errFrame, ok := frame.(frames.ErrorFrame); ok {
		log.Printf("Error running app: %v", errFrame.Error)
		if errFrame.Fatal {
			// Cancel all tasks downstream.
			s.PushFrame(frames.CancelFrame{}, processors.FrameDirectionDownstream)
			// Tell the task we should stop.
			s.upQueue <- frames.StopTaskFrame{}
		}
	}
}

// PipelineTask runs a pipeline.
type PipelineTask struct {
	ID       int
	Name     string
	pipeline processors.FrameProcessor
	params   PipelineParams
	finished bool
	downQueue chan frames.Frame
	upQueue   chan frames.Frame
	source   *TaskSource
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

var taskCounter int
var taskCounterMu sync.Mutex

func nextTaskID() int {
	taskCounterMu.Lock()
	defer taskCounterMu.Unlock()
	taskCounter++
	return taskCounter
}

func NewPipelineTask(pipeline processors.FrameProcessor, params PipelineParams) *PipelineTask {
	id := nextTaskID()
	ctx, cancel := context.WithCancel(context.Background())
	task := &PipelineTask{
		ID:        id,
		Name:      fmt.Sprintf("PipelineTask#%d", id),
		pipeline:  pipeline,
		params:    params,
		downQueue: make(chan frames.Frame, 128),
		upQueue:   make(chan frames.Frame, 128),
		ctx:       ctx,
		cancel:    cancel,
	}
	task.source = NewTaskSource(task.upQueue)
	task.source.Link(pipeline)
	return task
}

func (t *PipelineTask) HasFinished() bool {
	return t.finished
}

func (t *PipelineTask) StopWhenDone() {
	log.Printf("Task %s scheduled to stop when done", t.Name)
	t.QueueFrame(frames.EndFrame{})
}

func (t *PipelineTask) Cancel() {
	log.Printf("Canceling pipeline task %s", t.Name)
	t.source.PushFrame(frames.CancelFrame{}, processors.FrameDirectionDownstream)
	t.cancel()
}

func (t *PipelineTask) Run() {
	t.wg.Add(2)
	go t.processDownQueue()
	go t.processUpQueue()
	t.wg.Wait()
	t.finished = true
}

func (t *PipelineTask) QueueFrame(frame frames.Frame) {
	t.downQueue <- frame
}

func (t *PipelineTask) processDownQueue() {
	defer t.wg.Done()
	defer t.pipeline.Cleanup()
	defer t.cancel() // Signal other goroutines to stop when this one is done.
	defer close(t.downQueue)

	startFrame := frames.StartFrame{
		AllowInterruptions:     t.params.AllowInterruptions,
		EnableMetrics:          t.params.EnableMetrics,
		EnableUsageMetrics:     t.params.EnableUsageMetrics,
		ReportOnlyInitialTTFB:  t.params.ReportOnlyInitialTTFB,
	}
	t.source.ProcessFrame(startFrame, processors.FrameDirectionDownstream)

	running := true
	for running {
		select {
		case <-t.ctx.Done():
			return
		case frame, ok := <-t.downQueue:
			if !ok {
				running = false
				break
			}
			t.source.ProcessFrame(frame, processors.FrameDirectionDownstream)
			switch frame.(type) {
			case frames.StopTaskFrame, frames.EndFrame:
				running = false
			}
		}
	}
}

func (t *PipelineTask) processUpQueue() {
	defer t.wg.Done()
	defer close(t.upQueue)
	for {
		select {
		case <-t.ctx.Done():
			return
		case frame, ok := <-t.upQueue:
			if !ok {
				return
			}
			if _, ok := frame.(frames.StopTaskFrame); ok {
				t.QueueFrame(frames.StopTaskFrame{})
			}
		}
	}
}
