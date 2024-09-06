package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

func mapSecretsToViperConfigs() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func LoadConfigs() {
	//mapping secrets to viper configs
	mapSecretsToViperConfigs()
	env := viper.GetString("app.env")

	viper.SetConfigType("json")
	viper.AddConfigPath("resources/configs")
	if env == "" {
		log.Fatal().Msg("received empty env to load config for.")
	}
	viper.SetConfigName(env)
	err := viper.MergeInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read configs")
	}

}
