package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	log1 "github.com/sniperCore/core/log"
	"log"
)

/**
 * init es client
 */
func NewEs(conf *BaseConfig) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		//基础信息设置
		elastic.SetURL(conf.Host...),
		elastic.SetScheme(conf.Schema),
		elastic.SetBasicAuth(conf.BaseAuth.UserName, conf.BaseAuth.PassWord),

		//Sniff相关设置-客户端嗅探集群的状态
		elastic.SetSniff(conf.Sniff.IsSniff),
		elastic.SetSnifferTimeoutStartup(conf.Sniff.SnifferTimeoutStartup),
		elastic.SetSnifferTimeout(conf.Sniff.SnifferTimeout),
		elastic.SetSnifferInterval(conf.Sniff.SnifferInterval),

		//HealthCheck相关设置 - 健康检查相关设置
		elastic.SetHealthcheck(conf.HealthCheck.IsHealthCheck),
		elastic.SetHealthcheckTimeoutStartup(conf.HealthCheck.HealthcheckTimeoutStartup),
		elastic.SetHealthcheckTimeout(conf.HealthCheck.HealthcheckTimeout),
		elastic.SetInfoLog(log.New(log1.GetLogger(log1.SingletonMain).Writer(), "elastic ", log.LstdFlags)),
	)
	if err != nil {
		return &elastic.Client{}, err
	}

	for _, host := range conf.Host {
		_, _, err = client.Ping(host).Do(context.Background())
	}

	return client, err
}
