package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/vallox2mqtt/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vallox2mqtt",
		Short: "vallox2mqtt",
	}
)

func init() {
	cobra.OnInitialize(initialize)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatal().Err(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initialize() {
	if err := config.InitConfig(); err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize configuration")
	}
	viper.AutomaticEnv()

	initLogging()

	log.Debug().Str("LEVEL", zerolog.GlobalLevel().String()).Msg("Initialized logging")
	log.Debug().Msg("Initialized configuration")
}

func initLogging() {
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if viper.GetBool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
