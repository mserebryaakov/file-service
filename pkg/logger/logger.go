package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type MainLogHook struct{}

func (h *MainLogHook) Fire(entry *logrus.Entry) error {
	entry.Message = "Main: " + entry.Message
	return nil
}

func (h *MainLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func NewLogger(LogLvl string, hook logrus.Hook) *logrus.Entry {
	l := logrus.New()

	l.AddHook(hook)

	l.SetOutput(os.Stdout)

	l.SetFormatter(&logrus.TextFormatter{})

	logLevel, err := logrus.ParseLevel(LogLvl)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	l.SetLevel(logLevel)

	return logrus.NewEntry(l)
}
