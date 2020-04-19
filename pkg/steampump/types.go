package steampump

import (
	"sync"

	"github.com/m1cr0man/steampump/pkg/steam"
	"github.com/m1cr0man/steampump/pkg/steammesh"
)

type Config struct {
	Steam steam.Config     `json:"steam"`
	Mesh  steammesh.Config `json:"mesh"`
}

type TaskStatus int

const (
	TaskStatusQueued TaskStatus = iota
	TaskStatusRunning
	TaskStatusDone
)

func (s TaskStatus) String() string {
	return [3]string{"Queued", "Running", "Done"}[s]
}

type Task struct {
	Title    string
	Info     string
	Progress string
	Status   TaskStatus
	Gate     sync.Cond
}

func NewTask(title, info string) *Task {
	return &Task{
		Title: title,
		Info:  info,
		Gate:  sync.Cond{L: &sync.Mutex{}},
	}
}

const TaskQueueSize = 100

type WaitableTaskQueue struct {
	index  int
	outdex int
	tasks  [TaskQueueSize]*Task
	signal sync.Cond
}

func (q *WaitableTaskQueue) Put(task *Task) {
	// We want to preserve the tasks list so it can be queried
	// later. Only keeping 100
	q.tasks[q.index] = task
	q.index = (q.index + 1) % TaskQueueSize
	q.signal.Signal()
}

func (q *WaitableTaskQueue) Get() *Task {
	if q.index == q.outdex {
		q.signal.Wait()
	}
	outdex := q.outdex
	q.outdex = (q.outdex + 1) % TaskQueueSize
	return q.tasks[outdex]
}

func (q *WaitableTaskQueue) List() []*Task {
	return q.tasks[:]
}

func NewWaitableTaskQueue() *WaitableTaskQueue {
	return &WaitableTaskQueue{
		tasks:  [TaskQueueSize]*Task{},
		signal: sync.Cond{L: &sync.Mutex{}},
	}
}
