package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupLogger(config *Config) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
}
