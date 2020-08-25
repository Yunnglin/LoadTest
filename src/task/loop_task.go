package task

import (
	"LoadTest/src/interfaces"
	"LoadTest/src/util/queue"
	"sync"
	"time"
)

type LoopTask struct {
	// C is the concurrency level, the number of concurrent workers
	// to run.
	C int

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
			worker.(interfaces.Work).Close(l)
		} else {
			l.isStop = true
		}
	}
}

func (l *LoopTask) Wait() {
	l.wg.Wait()
}

func (l *LoopTask) runWorkers(manager interfaces.WorkManager) {
	l.wg.Add(l.C)
	for i := 0; i < l.C; i++ {
		work := manager.CreateWork()
		l.q.Push(work)
		go func() {
			work.Init(l)
			work.RunWorker(l)
			l.wg.Done()
		}()
	}
}
