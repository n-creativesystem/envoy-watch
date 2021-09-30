package main

import (
	"time"

	"github.com/n-creativesystem/envoy-watch/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetLevel(logrus.DebugLevel)
	cmd.Execute()
}
