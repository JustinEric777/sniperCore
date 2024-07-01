package es

import (
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sniperCore/core/container"
	"github.com/sniperCore/core/helper"
	"sync"
)

const SingletonMain = "es"

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

	conf, ok := args[1].(*BaseConfig)
	if !ok {
		return errors.New("args[1] is not config.EsConfig")
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
func setSingleton(diName string, conf *BaseConfig) (ins *elastic.Client, err error) {
	ins, err = NewEs(conf)
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
		es := getSingleton(k, false)
		es.Stop()
	}
	return nil
}

//获取单例
func getSingleton(diName string, lazy bool) *elastic.Client {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*elastic.Client)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(*BaseConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("es di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("es di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}

//外部通过注入别名获取资源，解耦资源的关系
func GetES(args ...string) *elastic.Client {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
