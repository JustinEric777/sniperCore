package log

import (
	"os"

	config2 "github.com/sniperCore/core/config"

	"github.com/sirupsen/logrus"

	hook2 "github.com/sniperCore/core/log/hook"
)

const DefaultPrefix = "access"

func InitLogger(config *LoggerConfig) (logger *logrus.Logger, err error) {
	logger = logrus.New()
	var prefix, level string
	var hook logrus.Hook
	dir := config2.GetAppPath() + "/"
	defaults := config.Default
	for _, logType := range defaults {
		switch logType {
		case "daily":
			//配置设置
			if config.Channels.Daily.Prefix == "" {
				prefix = DefaultPrefix
			}
			level = config.Channels.Daily.Level
			dir = dir + config.Channels.Daily.Path
			formatter := config.Channels.Daily.Format
			maxAge := config.Channels.Daily.Days
			rotationTime := config.Channels.Daily.RotationTime

			hook = hook2.NewRollHook(logger, dir, prefix, formatter, rotationTime, maxAge)
		case "single":
			//配置设置
			if config.Channels.Single.Prefix == "" {
				prefix = DefaultPrefix
			}
			level = config.Channels.Daily.Level
			level = config.Channels.Daily.Level
			dir = dir + config.Channels.Daily.Path
			formatter := config.Channels.Daily.Format
			maxAge := config.Channels.Daily.Days

			hook, err = hook2.NewLfsHook(dir, prefix, formatter, maxAge)
			if err != nil {
				return logger, err
			}
		}
	}

	//设置日志级别
	logLevel, err := logrus.ParseLevel(level)
	if err == nil {
		logger.SetLevel(logLevel)
	}

	//debug日志设置控制台输出
	if level == "debug" {
		logrus.SetOutput(os.Stdout)
	}

	//加载对应的hook
	logger.Hooks.Add(hook)

	return logger, nil
}
