package log

import (
	"github.com/sniperCore/core/config"

	"github.com/spf13/viper"
)

type LoggerConfig struct {
	Default  []string
	Channels Channels
}

type Channels struct {
	Single SingleLogConfig
	Daily  DailyLogConfig
}

//single log setting
type SingleLogConfig struct {
	Driver   string
	Level    string
	Format   string
	Prefix   string
	Path     string
	LinkPath string
	Days     int
}

//daily log setting
type DailyLogConfig struct {
	Driver       string
	Level        string
	Format       string
	Path         string
	Prefix       string
	LinkPath     string
	Days         int
	RotationTime int
}

/**
* log配置文件初始化
 */
func InitConfig() (*LoggerConfig, error) {
	var Config *LoggerConfig
	err := config.Conf.GetLocal("log")
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return nil, err
	}

	return Config, nil
}
