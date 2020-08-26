package task

import (
	"LoadTest/src/interfaces"
	"LoadTest/src/util/queue"
	"sync"
	"time"
)

type LoopTask struct {
	// Workers is the concurrency level, the number of concurrent workers
	// to run.
	Workers int

	// StartTime the time when start task
	StartTime time.Time

	//q the queue to store tasks
	q *queue.Queue

	//wg
	wg sync.WaitGroup

	//UserData
	UserData interface{}

	//IsStop
	isStop bool
}

func (l *LoopTask) IsStop() bool {
	return l.isStop
}

func (l *LoopTask) Run(manager interfaces.WorkManager) {
	l.StartTime = time.Now()
	l.q = queue.NewQueue()
	l.isStop = false
	l.runWorkers(manager)

}

func (l *LoopTask) Stop() {
	for !l.isStop {
		worker := l.q.Pop()
		if worker != nil {
			worker.(interfaces.Work).Close()
		} else {
			l.isStop = true
		}
	}
}

func (l *LoopTask) Wait() {
	l.wg.Wait()
}

func (l *LoopTask) runWorkers(manager interfaces.WorkManager) {
	l.wg.Add(l.Workers)
	for i := 0; i < l.Workers; i++ {
		work := manager.CreateWork()
		l.q.Push(work)
		go func() {
			work.Init()
			work.RunWorker()
			l.wg.Done()
		}()
	}
}
