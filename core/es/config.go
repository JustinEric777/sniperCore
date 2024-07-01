package es

import (
	"github.com/sniperCore/core/config"
	"time"
)

type ElasticConfigs struct {
	Connections map[string]*BaseConfig
}

type BaseConfig struct {
	Host        []string
	Schema      string
	MaxRetries  int
	BaseAuth    BaseAuth
	Sniff       Sniff
	HealthCheck HealthCheck
}

type BaseAuth struct {
	UserName string
	PassWord string
}

type Sniff struct {
	IsSniff               bool
	SnifferTimeoutStartup time.Duration
	SnifferTimeout        time.Duration
	SnifferInterval       time.Duration
}

type HealthCheck struct {
	IsHealthCheck             bool
	HealthcheckTimeoutStartup time.Duration
	HealthcheckTimeout        time.Duration
	HealthcheckInterval       time.Duration
}

/**
 * es配置文件初始化
 */
func InitConfig() (*ElasticConfigs, error) {
	var Configs *ElasticConfigs
	configure := config.Conf
	err := configure.Get("database")
	if err != nil {
		return nil, err
	}

	err = configure.UnmarshalKey("es", &Configs)
	if err != nil {
		return nil, err
	}

	return Configs, nil
}
