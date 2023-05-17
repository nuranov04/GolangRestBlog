package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.mod/pkg/logging"
	"sync"
)

type Config struct {
	Storage struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `yaml:"storage"`
	JWT struct {
		Secret string `yaml:"secret" env-required:"true"`
	}
	isDebug *bool `yaml:"is_debug env-required:true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIp string `yaml:"bind_ip" env-default:"0.0.0.0"`
		Port   string `yaml:"port" env-default:"8000"`
	} `yaml:"lister"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
