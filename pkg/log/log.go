package log

import (
	"flag"
	"io"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type LogInfo map[string]interface{}

func getLogger() *logrus.Logger {
	log := logrus.New()
	
	if flag.Lookup("test.v") != nil {
		log.SetOutput(io.Discard)
		return log
	}

	currDate := time.Now().Format("2006-01-02")

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "./data/logs/"+ currDate + ".log",
		logrus.WarnLevel: "./data/logs/"+ currDate + ".log",
		logrus.ErrorLevel: "./data/logs/"+ currDate + ".log",
	}

	log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return log 
}

func Info(fields LogInfo, info string) {
	var convFields map[string]interface{} = fields

	log := getLogger()

	log.WithFields(convFields).Info(info)
}

func Warn(fields LogInfo, info string) {
	var convFields map[string]interface{} = fields

	log := getLogger()

	log.WithFields(convFields).Warn(info)
}

func Fatal(fields LogInfo, info string) {
	var convFields map[string]interface{} = fields

	log := getLogger()

	log.WithFields(convFields).Fatal(info)
}

func Error(fields LogInfo, info string) {
	var convFields map[string]interface{} = fields

	log := getLogger()

	log.WithFields(convFields).Error(info)
}
