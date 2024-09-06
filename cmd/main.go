package main

import (
	"github.com/rs/zerolog/log"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/app"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/config"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/utils"
)

func main() {
	config.LoadConfigs()
	utils.InitLogger()
	application, err := app.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("error while initialising app")
	}
	if err := application.Start(); err != nil {
		log.Fatal().Err(err).Msg("error while starting app")
	}
}
