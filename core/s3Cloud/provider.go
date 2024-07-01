package s3Cloud

import (
	"errors"
	"fmt"
	"github.com/sniperCore/core/container"
	"github.com/sniperCore/core/helper"
	"sync"
)

const SingletonMain = "s3Cloud"

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

//注入logger日志类
func (p *provider) Register(args ...interface{}) error {
	diName, lazy, err := helper.TransformArgs(args...)
	if err != nil {
		return err
	}

	conf, ok := args[1].(*CloudStorageConfig)
	if !ok {
		return errors.New("args[1] is not config.s3CloudConfig")
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

func setSingleton(diName string, conf *CloudStorageConfig) (*S3Cloud, error) {
	consul, err := NewS3Cloud(conf)
	if err == nil {
		container.App.SetSingleton(diName, consul)
	}

	return consul, err
}

func (p *provider) Provides() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return helper.MapToArray(p.mp)
}

func getSingleton(diName string, lazy bool) *S3Cloud {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*S3Cloud)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(*CloudStorageConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("s3cloud di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("s3cloud di_name:%s err:%s", diName, err.Error()))
	}

	return ins
}

//外部通过注入别名获取资源，解耦资源的关系
func GetS3Cloud(args ...string) *S3Cloud {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
