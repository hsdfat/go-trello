package logger

// Reference: https://gist.github.com/gonzaloserrano/d360c1b195f7313b6d19

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"runtime"
	"strings"
)

var logger = logrus.New()

func Errorln(v ...any) {
	log.Println("[Error]: ", v)
}

func Debugln(v ...any) {
	log.Println("[Debug]: ", v)
}

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

func Fatalln(v ...any) {
	log.Fatalln("[Error]: ", v)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
