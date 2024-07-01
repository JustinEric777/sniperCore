package server

import (
	"github.com/sniperCore/core/config"
	"github.com/sniperCore/core/consul"
	"time"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Basic  *Basic
	Http   *HttpConfig
	Consul *consul.ConsulConfig
}

type HttpConfig struct {
	Addr         string
	Port         string
	Debug        bool
	IdleTimeout  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

type Basic struct {
	ServiceId   string
	ServiceName string
}

func InitConfig() (*ServiceConfig, error) {
	var Config *ServiceConfig
	err := config.Conf.GetLocal("server")
	if err != nil {
		return nil, err
	}
	err = viper.UnmarshalKey("servers", &Config)
	if err != nil {
		return nil, err
	}

	return Config, err
}
