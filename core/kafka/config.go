package kafka

import (
	"github.com/sniperCore/core/config"
)

type KafkaConfigs struct {
	Connections map[string]*KafkaConfig
}

type KafkaConfig struct {
	Hosts        []string `yaml:"hosts"`
	Topic        string   `yaml:"topic"`
	GroupId      string   `yaml:"groupId"`
	DialTimeout  int      `yaml:"dialTimeout"`
	ReadTimeout  int      `yaml:"readTimeout"`
	WriteTimeout int      `yaml:"writeTimeout"`
}

func InitConfig() (*KafkaConfigs, error) {
	var Configs *KafkaConfigs
	configure := config.Conf
	err := configure.Get("mq")
	if err != nil {
		return nil, err
	}

	err = configure.UnmarshalKey("kafka", &Configs)
	if err != nil {
		return nil, err
	}

	return Configs, nil
}
