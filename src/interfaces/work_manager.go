package interfaces

import "LoadTest/src/util/config"

type WorkManager interface {
	CreateWorker(workerId string) Work
	Finish(task Task)
	GetTask() *config.Task
}
