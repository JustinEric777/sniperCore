package hook

//根据日志级别进行文件切分
import (
	"fmt"
	"github.com/sniperCore/core/log/formatter"
	"io"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func NewLfsHook(logDir, name, format string, maxAge int) (*lfshook.LfsHook, error) {
	var (
		err         error
		infoWriter  io.Writer
		warnWriter  io.Writer
		errorWriter io.Writer
	)

	infoPath := fmt.Sprintf("%s/%s.%s", logDir, name, "INFO.%Y%m%d.log")
	linkInfoPath := fmt.Sprintf("%s/%s.%s", logDir, name, "INFO.log")
	infoWriter, err = rotatelogs.New(
		infoPath,
		rotatelogs.WithLinkName(linkInfoPath),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		rotatelogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	warnPath := fmt.Sprintf("%s/%s.%s", logDir, name, "WARN.%Y%m%d.log")
	linkWarnPath := fmt.Sprintf("%s/%s.%s", logDir, name, "WARN.log")
	warnWriter, err = rotatelogs.New(
		warnPath,
		rotatelogs.WithLinkName(linkWarnPath),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		rotatelogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	errorPath := fmt.Sprintf("%s/%s.%s", logDir, name, "ERROR.%Y%m%d.log")
	linkErrorPath := fmt.Sprintf("%s/%s.%s", logDir, name, "ERROR.log")
	errorWriter, err = rotatelogs.New(
		errorPath,
		rotatelogs.WithLinkName(linkErrorPath),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		rotatelogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	writerMap := lfshook.WriterMap{}
	if infoWriter != nil {
		writerMap[logrus.TraceLevel] = infoWriter
		writerMap[logrus.DebugLevel] = infoWriter
		writerMap[logrus.InfoLevel] = infoWriter
	}
	if warnWriter != nil {
		writerMap[logrus.WarnLevel] = warnWriter
	}
	if errorWriter != nil {
		writerMap[logrus.ErrorLevel] = errorWriter
		writerMap[logrus.FatalLevel] = errorWriter
		writerMap[logrus.PanicLevel] = errorWriter
	}

	if format == "json" {
		return lfshook.NewHook(writerMap, &formatter.JsonSpaceFormatter{}), nil
	} else {
		return lfshook.NewHook(writerMap, &formatter.TextSpaceFormatter{}), nil
	}
}
