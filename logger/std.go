package logger

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

type StdLogger interface {
	logrus.FieldLogger
	io.Writer
}

type stdLogger struct {
	*logrus.Logger
}

func (log *stdLogger) Write(p []byte) (n int, err error) {
	log.Logger.Print(string(p))
	return len(p), nil
}

func NewStdLogger(out io.Writer) StdLogger {
	log := logrus.New()
	log.SetOutput(out)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	return &stdLogger{
		Logger: log,
	}
}
