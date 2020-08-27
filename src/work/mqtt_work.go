package work

import (
	"LoadTest/src/util/log"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"sync"
)

type MqttWork struct {
	client       MQTT.Client
	lock         *sync.Mutex
	waitingQueue map[string]func(client MQTT.Client, msg MQTT.Message)
}


func (m *MqttWork) Connect(opts *MQTT.ClientOptions) error {
	log.Info.Println("Connecting")
	m.client = MQTT.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttWork) GetClient() MQTT.Client {
	return m.client
}

func (m *MqttWork) Finish() {
	m.client.Disconnect(250)
}

//Request will send a request to the server with the given
//	topic: the destination
//	body: message in byte
func (m *MqttWork) Request(topic string, body []byte) (MQTT.Message, error) {
	result := make(chan MQTT.Message)
	m.On(topic, func(client MQTT.Client, msg MQTT.Message) {
		result <- msg
	})
	m.GetClient().Publish(topic, 0, false, body)
	msg, ok := <-result
	if !ok {
		return nil, fmt.Errorf("client closed")
	}
	return msg, nil
}

// RequestNR request without a result
func (m *MqttWork) RequestNR(topic string, qos byte, retained bool, message string) {
	m.GetClient().Publish(topic, qos, retained, message)
}

func (m *MqttWork) On(topic string, f func(client MQTT.Client, msg MQTT.Message)) {
	m.lock.Lock()
	m.waitingQueue[topic] = f
	m.lock.Unlock()
}
