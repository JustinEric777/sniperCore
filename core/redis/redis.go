package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

//init redis client
func NewRedis(config *RedisBaseConfig) (*redis.Client, error) {
	err := optionCheck(config)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Database,
		PoolSize:     config.Options.PoolSize,
		DialTimeout:  config.Options.DialTimeout * time.Second,
		ReadTimeout:  config.Options.ReadTimeout * time.Second,
		WriteTimeout: config.Options.WriteTimeout * time.Second,
		PoolTimeout:  config.Options.PoolTimeout * time.Second,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

//redis相关配置处理
func optionCheck(config *RedisBaseConfig) error {
	if config.Host == "" {
		return errors.New("redis host is not empty")
	}
	if config.Options.PoolSize == 0 {
		config.Options.PoolSize = 1000
	}
	if config.Options.DialTimeout == 0 {
		config.Options.PoolSize = 30
	}
	if config.Options.ReadTimeout == 0 {
		config.Options.PoolSize = 30
	}
	if config.Options.WriteTimeout == 0 {
		config.Options.PoolSize = 30
	}

	return nil
}
