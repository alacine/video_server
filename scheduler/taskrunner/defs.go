package taskrunner

const (
	// ReadyToDispatch 可分配信号
	ReadyToDispatch = "d"

	// ReadyToExecute 可执行信号
	ReadyToExecute = "e"

	// Close 关闭信号
	Close = "c"

	// VideoDir 视频存放路径
	VideoDir = "./videos/"
)

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error
