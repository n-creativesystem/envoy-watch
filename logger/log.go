package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetLogLevel() logrus.Level {
	level, err := logrus.ParseLevel(viper.GetString("logLevel"))
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}
