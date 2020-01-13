package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLogger() {
	log.SetFormatter(&logrus.JSONFormatter{})
	// log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

func GetLogger() *logrus.Logger {
	return log
}
