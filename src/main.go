package main

import (
	"LoadTest/src/mqtt_task"
	"LoadTest/src/task"
	"LoadTest/src/util/config"
	"LoadTest/src/util/log"
	"LoadTest/src/util/queue"
	"LoadTest/src/util/result"
	"os"
	"os/signal"
)

func main() {
	tasks := config.LoadConfig("config/task_config.toml")
	result.GetResult()
	q := queue.NewQueue()
	c := make(chan os.Signal, 1)

	for i := range tasks.Tasks {
		manager := mqtt_task.NewManager(tasks.Tasks[i])
		loopTask := &task.LoopTask{
			Manager: manager,
		}
		go loopTask.Run()

		q.Push(loopTask)
	}

	go result.Run()

	log.Info.Println("Start load test")

	// 监听中断和终止信号
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	<-c
	log.Info.Println("Stopping")

	for true {
		item := q.Pop()
		t, ok := item.(task.LoopTask)
		if ok {
			t.Stop()
			t.Wait()
		}else{
			break
		}
	}
	result.Stop()
	log.Info.Println("Finished")
}
