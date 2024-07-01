package redis

import (
	"github.com/sniperCore/core/config"
	"time"
)

type RedisConfigs struct {
	Connections map[string]*RedisBaseConfig
}

type RedisBaseConfig struct {
	Host     string
	Port     string
	Password string
	Database int
	Options  redisOptionConfig
}

type redisOptionConfig struct {
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

func InitConfig() (*RedisConfigs, error) {
	var Configs *RedisConfigs

	configure := config.Conf
	err := configure.Get("database")
	if err != nil {
		return nil, err
	}

	err = configure.UnmarshalKey("redis", &Configs)
	return Configs, err
}
