package work

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"sync"
)

type MqttWork struct {
	client       MQTT.Client
	curId        int64
	lock         *sync.Mutex
	waitingQueue map[string]func(client MQTT.Client, msg MQTT.Message)
}

func (m *MqttWork) GetDefaultOptions(brokerAddr string) *MQTT.ClientOptions {
	m.curId = 0
	m.lock = &sync.Mutex{}
	m.waitingQueue = make(map[string]func(client MQTT.Client, msg MQTT.Message))
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerAddr)
	opts.SetClientID("1")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)
	opts.SetProtocolVersion(3)
	opts.SetAutoReconnect(false)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		m.lock.Lock()
		if callback, ok := m.waitingQueue[msg.Topic()]; ok {
			_, err := url.Parse(msg.Topic())
			if err != nil {
				fmt.Println("There is something wrong in parsing url: " + err.Error())
				delete(m.waitingQueue, msg.Topic())
			}
			go callback(client, msg)
		}
		m.lock.Unlock()
	})

	return opts
}
