package hook

import (
	"github.com/sniperCore/core/log/formatter"
	"time"

	"github.com/sirupsen/logrus"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type RollHook struct {
	logger       *logrus.Logger
	dir          string
	name         string
	format       string
	rotationTime int
	maxAge       int
}

//根据时间维度进行日志文件切分
func NewRollHook(logger *logrus.Logger, logDir, name, format string, rotationTime, maxAge int) *RollHook {
	rh := new(RollHook)
	rh.logger = logger
	rh.dir = logDir
	rh.name = name
	rh.format = format

	//格式化文件名
	rh.SetRollType(rotationTime)
	return rh
}

func (rh *RollHook) roll() error {
	//设置log writer
	logWriter, err := rotatelogs.New(
		rh.dir+"/"+rh.name+".log",
		rotatelogs.WithRotationTime(time.Duration(rh.rotationTime)*time.Second),
		rotatelogs.WithMaxAge(time.Duration(rh.maxAge)*24*time.Hour),
	)
	if err != nil {
		return err
	}
	rh.logger.SetOutput(logWriter)

	//设置格式化输出
	if rh.format == "json" {
		rh.logger.SetFormatter(&formatter.JsonSpaceFormatter{})
	} else {
		rh.logger.SetFormatter(&formatter.TextSpaceFormatter{})
	}

	return nil
}

func (rh *RollHook) SetRollType(rotationTime int) {
	if rotationTime < 86400 {
		rh.rotationTime = 3600
		rh.name = rh.name + "_" + "%Y%m%d%H"
	} else {
		rh.rotationTime = 86400
		rh.name = rh.name + "_" + "%Y%m%d"
	}
}

func (rh *RollHook) Fire(entry *logrus.Entry) error {
	return rh.roll()
}

func (rh *RollHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}
}
