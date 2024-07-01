package kafka

import (
	"errors"
	"fmt"
	"github.com/sniperCore/core/container"
	"github.com/sniperCore/core/helper"
	"sync"
)

const SingletonMain = "kafka"

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

	conf, ok := args[1].(*KafkaConfig)
	if !ok {
		return errors.New("args[1] is not config.KafkaConfig")
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
func setSingleton(diName string, conf *KafkaConfig) (ins *Kafka, err error) {
	ins, err = NewKafka(conf)
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
		kafka := getSingleton(k, false)
		if kafka != nil {
			kafka.Producer.Close()
			kafka.ConsumerGroup.Close()
		}
	}
	return nil
}

//获取单例
func getSingleton(diName string, lazy bool) *Kafka {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*Kafka)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(*KafkaConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("kafka di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("kafka di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}

//外部通过注入别名获取资源，解耦资源的关系
func GetKafka(args ...string) *Kafka {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
