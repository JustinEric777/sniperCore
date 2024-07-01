package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sniperCore/core/container"
	"github.com/sniperCore/core/helper"
	"sync"
)

const SingletonMain = "redis"

var Pr *provider

func init() {
	Pr = new(provider)
	Pr.mp = make(map[string]interface{})
}

type provider struct {
	mu sync.RWMutex
	mp map[string]interface{}
	dn string
}

//DB服务注册
func (p *provider) Register(args ...interface{}) error {
	diName, lazy, err := helper.TransformArgs(args...)
	if err != nil {
		return err
	}

	conf, ok := args[1].(*RedisBaseConfig)
	if !ok {
		return errors.New("args[1] is not config.RedisConfig")
	}

	p.mu.Lock()
	p.mp[diName] = args[1]
	if len(p.mp) == 1 {
		p.dn = diName
	}
	p.mu.Unlock()

	if !lazy {
		_, err = setSingleton(diName, conf)
	}

	return nil
}

//注入单例
func setSingleton(diName string, conf *RedisBaseConfig) (ins *redis.Client, err error) {
	ins, err = NewRedis(conf)
	if err == nil {
		container.App.SetSingleton(diName, ins)
	}
	return
}

//打印出注册过的服务信息
func (p *provider) Provides() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return helper.MapToArray(p.mp)
}

//释放资源
func (p *provider) Close() error {
	arr := p.Provides()
	for _, k := range arr {
		redisClient := getSingleton(k, false)
		if redisClient != nil {
			err := redisClient.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//获取单例
func getSingleton(diName string, lazy bool) *redis.Client {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*redis.Client)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(*RedisBaseConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("redis di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("redis di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}

//外部通过注入别名获取资源，解耦资源的关系
func GetRedis(args ...string) *redis.Client {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
