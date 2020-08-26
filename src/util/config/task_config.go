package config

import (
	"LoadTest/src/util/log"
	"github.com/BurntSushi/toml"
	"path/filepath"
	"sync"
)

var (
	cfg  *taskConfig
	once sync.Once
)

type worker struct {
	Workers     int
	QPS         int  `toml:"quest_per_second"`
	OnlyConnect bool `toml:"only_connect"`
	CPS         int  `toml:"connection_per_second"`
}

type request struct {
	Topic    string
	Qos      byte
	Retained bool
	Message  string
}

type server struct {
	Servers      []string
	ClientId     string
	Username     string
	Password     string
	CleanSession bool `toml:"clean_session"`
}

type taskConfig struct {
	Worker  worker
	Server  server
	Request request
}

func LoadConfig(path string) *taskConfig {
	once.Do(func() {
		filePath, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		log.Info.Printf("parse toml file once. filePath: %s \n", filePath)

		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			panic(err)
		}
	})
	return cfg
}

func GetConfig() *taskConfig {
	return cfg
}
