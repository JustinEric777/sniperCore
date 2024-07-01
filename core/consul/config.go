package consul

import (
	"github.com/sniperCore/core/config"
	"strconv"

	"github.com/spf13/viper"
)

var ConConfig *ConsulConfig

type ConsulConfig struct {
	ServiceId   string            `json:"serviceId"`
	ServiceName string            `json:"service_name"`
	ServiceHost string            `json:"service_host"`
	ServicePort int               `json:"service_port"`
	WorkerId    string            `json:"worker_id"`
	WorkerType  string            `json:"worker_type"`
	WorkerCount int               `json:"worker_count"`
	WorkerMeta  map[string]string `json:"worker_meta"`
	Tag         []string          `json:"tag"`
	Enabled     bool              `json:"enabled"`
	Host        string            `json:"host"`
	Port        int               `json:"port"`
	Token       string            `json:"token"`
	Check       *Check            `json:"check"`
}

type Check struct {
	Enabled    bool `json:"enabled"`
	Timeout    int  `json:"timeout"`
	Interval   int  `json:"interval"`
	RemoveTime int  `json:"remove_time"`
}

func InitConsulConfig() (*ConsulConfig, error) {
	configure := config.Conf
	err := configure.GetLocal("server")
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("servers.basic", &ConConfig)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("servers.consul", &ConConfig)
	if err != nil {
		return nil, err
	}
	serviceHost, _ := configure.GetString("server.servers.http.host")
	if serviceHost != "" {
		ConConfig.ServiceHost = serviceHost
	}
	servicePort, _ := configure.GetString("server.servers.http.port")
	if servicePort != "" {
		port, _ := strconv.Atoi(servicePort)
		ConConfig.ServicePort = port
	}

	return ConConfig, nil
}
