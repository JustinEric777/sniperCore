package log

import (
	"errors"
	"fmt"
	"github.com/sniperCore/core/container"
	"github.com/sniperCore/core/helper"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const SingletonMain = "logger"

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

	conf, ok := args[1].(*LoggerConfig)
	if !ok {
		return errors.New("args[1] is not config.loggerConfig")
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

func setSingleton(diName string, conf *LoggerConfig) (*logrus.Logger, error) {
	logger, err := InitLogger(conf)
	if err == nil {
		container.App.SetSingleton(diName, logger)
	}

	return logger, err
}

func (p *provider) Provides() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return helper.MapToArray(p.mp)
}

func getSingleton(diName string, lazy bool) *logrus.Logger {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*logrus.Logger)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(*LoggerConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("logger di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("logger di_name:%s err:%s", diName, err.Error()))
	}

	return ins
}

//释放资源
func (p *provider) Close() error {
	arr := p.Provides()
	for _, k := range arr {
		logger := getSingleton(k, false)
		if logger != nil {
			log, ok := logger.Out.(*os.File)
			if ok {
				log.Sync()
				log.Close()
			}
		}
	}
	return nil
}

//外部通过注入别名获取资源，解耦资源的关系
func GetLogger(args ...string) *logrus.Logger {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
