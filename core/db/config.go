package db

import (
	"github.com/sniperCore/core/config"
)

type DbConfigs struct {
	Default     string
	Connections map[string]*DbConfig
}

type DbConfig struct {
	//驱动类型，目前支持mysql、postgres、mssql、sqlite3
	Driver       string         `json:"driver"`
	BaseConfig   DbBaseConfig   `mapstructure:"base" json:"base_config"`
	OptionConfig DbOptionConfig `mapstructure:"option" json:"option_config"`
}

type DbBaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

type DbOptionConfig struct {
	MaxIdle        int    `json:"max_idle"`
	MaxConns       int    `json:"max_conns"`
	IdleTimeout    int    `json:"idle_timeout"`
	ConnectTimeout int    `json:"connect_timeout"`
	Charset        string `json:"charset"`
}

func InitConfig() (*DbConfigs, error) {
	var Configs *DbConfigs
	configure := config.Conf
	err := configure.Get("database")
	if err != nil {
		return nil, err
	}
	err = configure.UnmarshalKey("database", &Configs)
	if err != nil {
		return nil, err
	}

	return Configs, nil
}
