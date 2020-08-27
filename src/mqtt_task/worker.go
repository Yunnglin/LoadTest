package mqtt_task

import (
	"LoadTest/src/util/log"
	"LoadTest/src/work"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"time"
)

type Worker struct {
	mqttWork *work.MqttWork
	manager  *Manager
	// requests per second
	QPS         int
	ConnectOnly bool
	closeSignal bool
	deviceId    string
}

func (w *Worker) Init() {
	w.closeSignal = false
}

func (w *Worker) RunWorker() {
	for !w.closeSignal {
		var throttle <-chan time.Time
		if w.QPS > 0 {
			throttle = time.Tick(time.Duration(1e6/(w.QPS)) * time.Microsecond)
		}

		if w.QPS > 0 {
			<-throttle
		}
		w.work()
	}
}

func (w *Worker) Close() {
	w.closeSignal = true
}

func (w *Worker) work() {
	if !w.ConnectOnly {
		request := w.manager.GetTask().Request
		request.Topic = strings.Replace(request.Topic, "{device_id}", w.deviceId, 1)
		w.mqttWork.RequestNR(request.Topic, request.Qos, request.Retained, request.Message)
	}

}

func NewWorker(manager *Manager, opts *MQTT.ClientOptions) *Worker {

	worker := new(Worker)
	worker.manager = manager
	worker.mqttWork = new(work.MqttWork)

	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		log.Error.Println("Connection Lost!", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		log.Info.Println("Connected!")
	})
	opts.SetTLSConfig(manager.Cert())
	err := worker.mqttWork.Connect(opts)
	if err != nil {
		log.Error.Println(err.Error())
	}
	return worker
}
