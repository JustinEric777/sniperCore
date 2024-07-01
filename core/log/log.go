package log

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

func Logf(level logrus.Level, format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Logf(level, format, args)
}

func Tracef(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Tracef(format, args)
}

func Debugf(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Debugf(format, args)
}

func Infof(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Infof(format, args)
}

func Printf(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Printf(format, args)
}

func Warnf(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Warnf(format, args)
}

func Warningf(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Warningf(format, args)
}

func Errorf(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Errorf(format, args)
}

func Fatalf(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Fatalf(format, args)
}

func Panicf(format string, args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Panicf(format, args)
}

func Trace(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Trace(args)
}

func Debug(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Debug(args)
}

func Info(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Info(args)
}

func Print(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Print(args)
}

func Warn(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Warn(args)
}

func Warning(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Warn(args...)
}

func Error(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Error(args)
}

func Fatal(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Fatal(args)
}

func Panic(args ...interface{}) {
	caller := getReportCaller()
	GetLogger(SingletonMain).WithFields(caller).Panic(args)
}

func getReportCaller() map[string]interface{} {
	caller := make(map[string]interface{})

	_, file, line, _ := runtime.Caller(2)
	caller["file"] = file
	caller["line"] = line

	return caller
}
