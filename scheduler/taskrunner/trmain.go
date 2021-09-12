package taskrunner

import (
	"time"
)

// Worker ...
type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

// NewWorker ...
func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {
	//for {
	//select {
	//case <-w.ticker.C:
	//go w.runner.StartAll()
	//}
	//}
	for range w.ticker.C {
		go w.runner.StartAll()
	}
}

// Start 新建 Worker 并启动
func Start() {
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)
	go w.startWorker()
}
