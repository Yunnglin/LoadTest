package interfaces

type Work interface {
	Init()
	RunWorker()
	Close()
}