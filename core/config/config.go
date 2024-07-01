package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

/**
 * 依托viper进行配置文件封装，同时兼容consul/etcd等
 */
var (
	Conf       *Config
	viperCloud *viper.Viper
)

type Config struct {
	IsCloud bool
	Type    string
	Host    string
	Port    int
	Prefix  string
}

/**
 * 初始化设置配置文件目录和格式
 */
func init() {
	dir := GetAppPath()
	viper.AddConfigPath(dir + "/config")
	viper.SetConfigType("yaml")
}

/**
 * 根据配置文件名和对应的key获取对应的值
 */
func (config *Config) Get(key string) error {
	err := config.isCloud()
	if err != nil {
		return err
	}

	if Conf.IsCloud {
		err = config.readRemoteConfig(key)
	} else {
		err = config.GetLocal(key)
	}
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) GetLocal(key string) error {
	file, _, err := getConfigFileName(key)
	if err != nil {
		return err
	}

	viper.SetConfigName(file)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

// 判断是否走配置中心
func (config *Config) isCloud() error {
	err := config.GetLocal("server")
	if err != nil {
		return err
	}
	err = viper.UnmarshalKey("config", &Conf)
	if err != nil {
		return err
	}

	return nil
}

// 读取配置中心相应的配置信息
func (config *Config) readRemoteConfig(key string) error {
	if Conf.IsCloud {
		viperCloud = viper.New()
		cloudType, addr, keyPath := Conf.Type, fmt.Sprintf("%s:%d", Conf.Host, Conf.Port), Conf.Prefix+key
		valid := cloudTypeValid(cloudType)
		if !valid {
			return errors.New("this config center is not supported")
		}
		err := viperCloud.AddRemoteProvider(cloudType, addr, keyPath)
		if err != nil {
			return err
		}
		viperCloud.SetConfigType("yaml")
		if err != nil {
			return err
		}
		err = viperCloud.ReadRemoteConfig()
		if err != nil {
			// 空文件跳过
			if err.Error() == "Remote Configurations Error: No Files Found" {
				return nil
			}
			return err
		}
	}

	return nil
}

func cloudTypeValid(cloudType string) bool {
	switch cloudType {
	case "etcd", "etcd3", "consul":
		return true
	}

	return false
}

func (config *Config) UnmarshalKey(key string, rawVal interface{}) error {
	if Conf.IsCloud {
		err := viperCloud.UnmarshalKey(key, rawVal)
		if err != nil {
			return err
		}
	}
	err := viper.UnmarshalKey(key, rawVal)
	if err != nil {
		return err
	}

	return nil
}

/**
 * 根据配置文件名和对应的key获取对应的值
 */
func (config *Config) GetString(key string) (string, error) {
	file, configKey, err := getConfigFileName(key)
	if err != nil {
		return "", err
	}
	viper.SetConfigName(file)
	err = viper.ReadInConfig()
	if err != nil {
		return "", err
	}

	value := viper.GetString(configKey)
	return value, nil
}

/**
 * 切割获取文件名和key
 */
func getConfigFileName(key string) (string, string, error) {
	fileArr := strings.Split(key, ".")
	if len(fileArr) == 0 {
		return "", "", errors.New("config file is ")
	}

	keyArr := fileArr[1:]
	configKey := strings.Join(keyArr, ".")
	return fileArr[0], configKey, nil
}

/**
 * 返回对应key的map映射
 */
func (config *Config) GetStringMap(key string) (map[string]interface{}, error) {
	keyMap := make(map[string]interface{})
	file, configKey, err := getConfigFileName(key)
	if err != nil {
		return keyMap, err
	}
	viper.SetConfigName(file)
	err = viper.ReadInConfig()
	if err != nil {
		return keyMap, err
	}

	return viper.GetStringMap(configKey), nil
}

func GetAppPath() string {
	var appPath string

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	appPath = path[:index]

	//linux
	if strings.HasPrefix(appPath, "/tmp/go-build") {
		appPath, _ = os.Getwd()
	}

	//mac
	if strings.HasPrefix(appPath, "/private/var/folders/") {
		appPath, _ = os.Getwd()
	}

	return appPath
}
