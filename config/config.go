package config

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("ha.discovery_prefix", "homeassistant")
	viper.SetDefault("ha.device_id", "vallox_digit_se")
	viper.SetDefault("vallox.port", "/dev/ttyUSB0")
	viper.SetDefault("vallox.baudRate", 9600)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	return nil
}
