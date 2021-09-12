package taskrunner

// Runner ...
type Runner struct {
	Controller controlChan // 生产者和消费者交换信息
	Error      controlChan
	Data       dataChan // 具体任务数据
	dataSize   int
	longLived  bool
	Dispatcher fn
	Executor   fn
}

// NewRunner ...
func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1), // 非阻塞带 buffer 的 channel
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		dataSize:   size,
		longLived:  longlived,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()
	for {
		select {
		case c := <-r.Controller:
			if c == ReadyToDispatch {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- Close
				} else {
					r.Controller <- ReadyToExecute
				}
			}
			if c == ReadyToExecute {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- Close
				} else {
					r.Controller <- ReadyToDispatch
				}
			}
		case e := <-r.Error:
			if e == Close {
				return
			}
		}
	}
}

// StartAll 启动 Runner
func (r *Runner) StartAll() {
	r.Controller <- ReadyToDispatch
	r.startDispatch()
}
