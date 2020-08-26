package main

import (
	"LoadTest/src/mqtt_task"
	"LoadTest/src/task"
	"LoadTest/src/util/config"
	"LoadTest/src/util/log"
	"os"
	"os/signal"
)

func main() {
	config.LoadConfig("config/task_config.toml")

	loopTask := task.LoopTask{
		Workers: config.GetConfig().Worker.Workers,
	}
	manager := mqtt_task.NewManager()
	log.Info.Println("Start load test")
	c := make(chan os.Signal, 1)
	go func() {
		loopTask.Run(manager)
	}()
	// 监听中断和终止信号
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	<-c
	log.Info.Println("Stopping")
	loopTask.Stop()
	loopTask.Wait()
	log.Info.Println("Finished")
}
