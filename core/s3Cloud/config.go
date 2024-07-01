package s3Cloud

import (
	"fmt"
	"github.com/sniperCore/core/config"
)

type S3CloudConfigs struct {
	Disks map[string]*CloudStorageConfig `json:"disks"`
}

type CloudStorageConfig struct {
	Driver   string `json:"driver"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Token    string `json:"token"`
	Region   string `json:"region"`
	Url      string `json:"url"`
	EndPoint string `json:"end_point" mapstructure:"end_point"`
	Bucket   string `json:"bucket"`
}

func InitConfig() (*S3CloudConfigs, error) {
	var Configs *S3CloudConfigs
	configure := config.Conf
	err := configure.Get("s3Cloud")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = configure.UnmarshalKey("s3Cloud", &Configs)
	if err != nil {
		return nil, err
	}

	return Configs, nil
}
