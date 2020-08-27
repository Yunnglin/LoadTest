package task

import (
	"LoadTest/src/interfaces"
	"LoadTest/src/util/log"
	"LoadTest/src/util/queue"
	"strconv"
	"sync"
	"time"
)

type LoopTask struct {
	Manager interfaces.WorkManager

	// StartTime the time when start task
	StartTime time.Time

	//q the queue to store tasks
	q *queue.Queue

	//wg
	wg sync.WaitGroup

	//IsStop
	isStop bool
}

func (l *LoopTask) IsStop() bool {
	return l.isStop
}

func (l *LoopTask) Run() {
	l.StartTime = time.Now()
	l.q = queue.NewQueue()
	l.isStop = false
	l.runWorkers()

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

func (l *LoopTask) runWorkers() {
	workers := l.Manager.GetTask().Worker.Workers
	startId := l.Manager.GetTask().Device.DeviceId[0]
	endId := l.Manager.GetTask().Device.DeviceId[1]

	l.wg.Add(workers)

	for i := 0; i < workers; i++ {
		if startId+i > endId {
			log.Warning.Println("Worker numbers exceed device id numbers")
			continue
		}
		worker := l.Manager.CreateWorker(strconv.Itoa(startId + i))
		l.q.Push(worker)
		go func() {
			worker.Init()
			worker.RunWorker()
			l.wg.Done()
		}()
	}
}
