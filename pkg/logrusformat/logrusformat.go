package logrusformat

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct {
	logrus.Formatter

	TimestampFormat string
}

func (f LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s]\t %s: %s\n",
		strings.ToUpper(entry.Level.String()), entry.Time.Format(f.TimestampFormat), entry.Message)), nil
}
