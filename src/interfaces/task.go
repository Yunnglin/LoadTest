package interfaces

type Task interface {
	Run(manager WorkManager)
	Stop()
	Wait()
}