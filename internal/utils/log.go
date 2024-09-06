package utils

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitLogger() {
	logLevel, err := zerolog.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		logLevel = zerolog.DebugLevel
	}
	//set log level
	log.Logger = log.Logger.Level(logLevel)
	//global contexts
	if len(viper.GetString("app.name")) > 0 {
		log.Logger = log.Logger.With().
			Str("service", viper.GetString("app.env")).
			Logger()
	}

	zerolog.DefaultContextLogger = &log.Logger
}
