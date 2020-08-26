package work

import (
	"LoadTest/src/util/config"
	"LoadTest/src/util/log"
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
	opts         *MQTT.ClientOptions
}

func (m *MqttWork) GetDefaultOptions() *MQTT.ClientOptions {
	m.curId = 0
	m.lock = &sync.Mutex{}
	m.waitingQueue = make(map[string]func(client MQTT.Client, msg MQTT.Message))
	server := config.GetConfig().Server
	opts := MQTT.NewClientOptions()
	opts.AddBroker(server.Servers[0])
	//opts.SetClientID(server.ClientId)
	opts.SetUsername(server.Username)
	opts.SetPassword(server.Password)
	opts.SetCleanSession(server.CleanSession)
	opts.SetProtocolVersion(3)
	opts.SetAutoReconnect(false)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		m.lock.Lock()
		if callback, ok := m.waitingQueue[msg.Topic()]; ok {
			_, err := url.Parse(msg.Topic())
			if err != nil {
				log.Error.Println("There is something wrong in parsing url: " + err.Error())
				delete(m.waitingQueue, msg.Topic())
			}
			go callback(client, msg)
		}
		m.lock.Unlock()
	})

	return opts
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
	m.curId += 1
	topic = fmt.Sprintf("%s/%d", topic, m.curId)
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
	m.curId += 1
	topic = fmt.Sprintf("%s/%d", topic, m.curId)
	m.GetClient().Publish(topic, qos, retained, message)
}

func (m *MqttWork) On(topic string, f func(client MQTT.Client, msg MQTT.Message)) {
	m.lock.Lock()
	m.waitingQueue[topic] = f
	m.lock.Unlock()
}
