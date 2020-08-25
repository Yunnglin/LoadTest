package interfaces

type Work interface {
	Init(task Task)
	RunWorker(task Task)
	Close(task Task)
}