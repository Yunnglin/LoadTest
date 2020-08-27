package result

import (
	"LoadTest/src/util/log"
	"sync"
	"time"
)

type Result struct {
	ConnectionCount int
	PublishingCount int
}

var (
	result *Result
	once   sync.Once
	lock   sync.Mutex
	isStop bool
)

func GetResult() *Result {
	once.Do(func() {
		isStop = false
		result = &Result{
			ConnectionCount: 0,
			PublishingCount: 0,
		}
	})
	return result
}

func AddConnection(count int) {
	lock.Lock()
	result.ConnectionCount += count
	lock.Unlock()
}

func AddPublishing(count int) {
	lock.Lock()
	result.PublishingCount += count
	lock.Unlock()
}

func PrintResult() {
	log.Info.Printf("Connection Count: %d", result.ConnectionCount)
	log.Info.Printf("Publishing Count: %d", result.PublishingCount)
}

func Run(){
	for !isStop {
		time.Sleep(1*time.Second)
		PrintResult()
	}
}

func Stop(){
	isStop = true
}
