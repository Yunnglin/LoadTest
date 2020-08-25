package interfaces

type WorkManager interface {
	CreateWork() Work
	Finish(task Task)
}
