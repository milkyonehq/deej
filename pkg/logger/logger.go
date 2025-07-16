package logger

import log "github.com/sirupsen/logrus"

func Init(level string) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	logLevel, err := log.ParseLevel(level)
	if err != nil {
		logLevel = log.InfoLevel
		log.Warnf("Invalid log level: %s, using default: %s", level, logLevel.String())
	}

	log.SetLevel(logLevel)
}
