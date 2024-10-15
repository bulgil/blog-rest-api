package logger

import (
	"log"
	"os"
	"time"

	"github.com/bulgil/blog-rest-api/pkg/logrusformat"
	"github.com/sirupsen/logrus"
)

func NewLogger(env string) *logrus.Logger {
	var logger = &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrusformat.LogFormatter{
			TimestampFormat: time.TimeOnly,
		},
	}

	switch env {
	case "dev":
		logger.SetLevel(logrus.DebugLevel)
	case "prod":
		logger.SetLevel(logrus.InfoLevel)
	}

	log.Println("logger initialized")
	return logger
}
