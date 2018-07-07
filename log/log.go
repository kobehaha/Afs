package log

import (
	l4g "github.com/cfxks1989/log4go"
	"github.com/kobehaha/Afs/constant"
	"os"
	"path"
)

var logger l4g.Logger

func Init() {

	loglevelMap := map[string]l4g.Level{
		"DEBUG":   l4g.DEBUG,
		"TRACE":   l4g.TRACE,
		"INFO":    l4g.INFO,
		"WARNING": l4g.WARNING,
		"ERROR":   l4g.ERROR,
	}

	logDir := ""
	logLevel := ""
	if os.Getenv("LOG_DIR") == "" {
		logDir = constant.STORAGE_LOG_DIR
	} else {
		logDir = os.Getenv("LOG_DIR")
	}
	if os.Getenv("LOG_LEVEL") == "" {
		logLevel = constant.STORAGE_LOG_LEVEL
	} else {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	logFileName := path.Join(logDir, "afs.log")
	logger = make(l4g.Logger)
	fileLogWriter := l4g.NewFileLogWriter(logFileName, false)
	fileLogWriter.SetFormat("[%D %T] [%L] (%S) %M")
	fileLogWriter.SetRotate(true)
	fileLogWriter.SetRotateDaily(true)

	if loglevel, found := loglevelMap[logLevel]; found {
		logger.AddFilter("file", loglevel, fileLogWriter)
		if loglevel == l4g.DEBUG {
			logger.AddFilter("file", loglevel, fileLogWriter)
		}
	} else {
		logger.AddFilter("file", l4g.INFO, fileLogWriter)
	}

}

func GetLogger() l4g.Logger {
	return logger
}
