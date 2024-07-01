package bootstrap

import (
	"fmt"
	"github.com/sniperCore/core/consul"
	"github.com/sniperCore/core/container"
	"github.com/sniperCore/core/db"
	"github.com/sniperCore/core/es"
	"github.com/sniperCore/core/helper"
	"github.com/sniperCore/core/kafka"
	"github.com/sniperCore/core/log"
	"github.com/sniperCore/core/redis"
	"github.com/sniperCore/core/s3Cloud"
)

var App *container.Container

func Bootstrap() error {
	helper.Println("Application bootstrap start init ...")

	//初始化容器
	App = container.App

	//step 0 : 注入log相关服务
	err := registerLogger()
	if err != nil {
		return err
	}

	//step 1 : 注入注册中心
	err = registerConsul()
	if err != nil {
		return err
	}

	//注入数据库DB服务
	err = registerDB()
	if err != nil {
		return err
	}

	//注入es
	err = registerES()
	if err != nil {
		return err
	}

	//注入redis
	err = registerRedis()
	if err != nil {
		return err
	}

	//注入kafka
	err = registerKafka()
	if err != nil {
		return err
	}

	//注入s3Cloud
	err = registerS3Cloud()
	if err != nil {
		return err
	}

	helper.Println("Application bootstrap init success...")

	return nil
}

//日志
func registerLogger() error {
	loggerConfig, err := log.InitConfig()
	if err != nil {
		return fmt.Errorf("Logger configs init error, err = %d", err)
	}
	if loggerConfig != nil {
		err = log.Pr.Register(log.SingletonMain, loggerConfig)
		if err != nil {
			return fmt.Errorf("Logger service provider register error, err = %d", err)
		}
		helper.Println("[logger]: Application bootstrap logger register success ...")
	}

	return nil
}

// consul
func registerConsul() error {
	consulConfig, err := consul.InitConsulConfig()
	if err != nil {
		return fmt.Errorf("Consul configs init error, err = %d", err)
	}
	if consulConfig != nil {
		err = consul.Pr.Register(consul.SingletonMain, consulConfig)
		if err != nil {
			return fmt.Errorf("Consul service provider register error, err = %d", err)
		}
		helper.Println("[consul]: Application bootstrap consul register success ...")
	}

	return nil
}

//DB
func registerDB() error {
	dbConfigs, err := db.InitConfig()
	if err != nil {
		return fmt.Errorf("DB configs init error , err = %d", err)
	}
	if dbConfigs != nil && len(dbConfigs.Connections) != 0 {
		for dname, dbConfig := range dbConfigs.Connections {
			err = db.Pr.Register(dname, dbConfig)
			if err != nil {
				return fmt.Errorf("DB [%s] service provider register error, err = %d", dname, err)
			}
		}
		helper.Println("[db]: Application bootstrap DB register success ...")
	}

	return nil
}

//es
func registerES() error {
	esConfigs, err := es.InitConfig()
	if err != nil {
		return fmt.Errorf("ES configs init error , err = %d", err)
	}
	if esConfigs != nil && len(esConfigs.Connections) != 0 {
		for ename, esConfig := range esConfigs.Connections {
			err = es.Pr.Register(ename, esConfig)
			if err != nil {
				return fmt.Errorf("ES [%s] service provider register error, err = %d", ename, err)
			}
			helper.Println("[es]: Application bootstrap ES register success ...")
		}
	}

	return nil
}

//redis
func registerRedis() error {
	redisConfigs, err := redis.InitConfig()
	if err != nil {
		return fmt.Errorf("Redis configs init error , err = %d", err)
	}
	if redisConfigs != nil && len(redisConfigs.Connections) != 0 {
		for rname, reidsConfig := range redisConfigs.Connections {
			err = redis.Pr.Register(rname, reidsConfig)
			if err != nil {
				return fmt.Errorf("Redis [%s] service provider register error, err = %d", rname, err)
			}
		}
		helper.Println("[redis]: Application bootstrap Redis register success ...")
	}

	return nil
}

// kafka
func registerKafka() error {
	kafkaConfigs, err := kafka.InitConfig()
	if err != nil {
		return fmt.Errorf("Kafka configs init error , err = %d", err)

	}

	if kafkaConfigs != nil && len(kafkaConfigs.Connections) != 0 {
		for kname, kafkaConfig := range kafkaConfigs.Connections {
			err = kafka.Pr.Register(kname, kafkaConfig)
			if err != nil {
				return fmt.Errorf("Kafka [%s] service provider register error, err = %d", kname, err)
			}
		}
		helper.Println("[kafka]: Application bootstrap Kafka register success ...")
	}

	return nil
}

// s3cloud
func registerS3Cloud() error {
	s3CloudConfigs, err := s3Cloud.InitConfig()
	if err != nil {
		return fmt.Errorf("s3Cloud configs init error , err = %d", err)

	}

	if s3CloudConfigs != nil && len(s3CloudConfigs.Disks) != 0 {
		for sName, s3CloudConfig := range s3CloudConfigs.Disks {
			err = s3Cloud.Pr.Register(sName, s3CloudConfig)
			if err != nil {
				return fmt.Errorf("s3Cloud [%s] service provider register error, err = %d", sName, err)
			}
		}
		helper.Println("[s3Cloud]: Application bootstrap s3Cloud register success ...")
	}

	return nil
}
