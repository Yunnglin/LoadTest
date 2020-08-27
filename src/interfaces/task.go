package interfaces

type Task interface {
	Run()
	Stop()
	Wait()
}