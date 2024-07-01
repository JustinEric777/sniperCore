package formatter

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type TextSpaceFormatter struct {
}

/**
 * 自定义日志的格式化格式
 */
func (f *TextSpaceFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	time := entry.Time.Format("2006-01-02 15:04:05.000")
	level := "[" + strings.Title(entry.Level.String()) + "]"
	entry.Message = strings.Trim(entry.Message, "[")
	entry.Message = strings.Trim(entry.Message, "]")
	if _, ok := entry.Data["file"]; ok && entry.Data["line"].(int) > 0 {
		fileName := entry.Data["file"]
		line := "line:" + strconv.Itoa(entry.Data["line"].(int))
		b.WriteString(fmt.Sprintf("%s %s %s %v %v", time, level, fileName, line, entry.Message))
	} else {
		b.WriteString(fmt.Sprintf("%s %s %v", time, level, entry.Message))
	}

	b.WriteString("\n")
	return b.Bytes(), nil
}
