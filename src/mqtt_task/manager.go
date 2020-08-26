package mqtt_task

import (
	"LoadTest/src/interfaces"
	"crypto/tls"
	"io"
	"sync"
)

type Manager struct {
	Writer io.Writer
	cert   *tls.Config
	lock   sync.RWMutex
}

func NewManager() *Manager{
	manager := &Manager{}
	return manager
}

func (m *Manager) CreateWork() interfaces.Work {
	return NewWorker(m)
}

func (m *Manager) Finish(task interfaces.Task) {

}

func (m *Manager) Cert() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
}
