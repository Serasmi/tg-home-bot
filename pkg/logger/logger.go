package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

var defaultLogger = New()

type Entry = logrus.Entry
type Fields = logrus.Fields

type Logger struct {
	*logrus.Entry
}

func New() *Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(logrus.DebugLevel)

	return &Logger{logrus.NewEntry(l)}
}

func GetLogger() *Logger {
	return defaultLogger
}

func (l *Logger) SetLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		l.Errorf("unable to parse loglevel: %s", lvl)
		return
	}

	l.Logger.SetLevel(level)
}

func (l *Logger) SetOut(writer io.Writer) {
	l.Logger.Out = writer
}
