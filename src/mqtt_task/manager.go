package mqtt_task

import (
	"LoadTest/src/interfaces"
	"LoadTest/src/util/config"
	"LoadTest/src/util/log"
	"crypto/tls"
	JWT "github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io"
	"strings"
	"sync"
)

type Manager struct {
	Writer io.Writer
	cert   *tls.Config
	lock   sync.RWMutex
	task   *config.Task
}

func (m *Manager) GetTask() *config.Task {
	return m.task
}

func NewManager(task *config.Task) *Manager {
	manager := &Manager{}
	manager.task = task
	return manager
}

func (m *Manager) CreateWorker(deviceId string) interfaces.Work {
	deviceCfg := m.task.Device
	workerCfg := m.task.Worker
	serverCfg := m.task.Server
	//处理clientId
	clientId := serverCfg.ClientId
	clientId = strings.Replace(clientId, "{project_id}", deviceCfg.ProjectId, 1)
	clientId = strings.Replace(clientId, "{hub_id}", deviceCfg.HubId, 1)
	clientId = strings.Replace(clientId, "{device_id}", deviceId, 1)
	//使用JWT作为password
	token := JWT.NewWithClaims(JWT.SigningMethodHS256, JWT.MapClaims{
		"id": deviceId,
	})
	tokenString, err := token.SignedString([]byte(deviceCfg.Secret))
	if err != nil {
		log.Error.Println(err.Error())
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(serverCfg.Broker)
	opts.SetUsername(deviceId)
	opts.SetPassword(tokenString)
	opts.SetClientID(clientId)
	opts.SetCleanSession(serverCfg.CleanSession)
	opts.SetProtocolVersion(3)
	opts.SetAutoReconnect(false)

	worker := NewWorker(m, opts)
	worker.QPS = workerCfg.QPS
	worker.ConnectOnly = workerCfg.OnlyConnect
	worker.deviceId = deviceId
	return worker
}

func (m *Manager) Finish(task interfaces.Task) {

}

func (m *Manager) Cert() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
}
