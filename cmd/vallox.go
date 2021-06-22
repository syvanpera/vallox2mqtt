package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/vallox2mqtt/vallox"
	"go.bug.st/serial"
)

var (
	valloxCmd = &cobra.Command{
		Use:     "vallox",
		Aliases: []string{"v"},
		Short:   "Vallox commands",
	}

	valloxTestCmd = &cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Test Vallox connection",
		Args:    cobra.NoArgs,
		Run:     runValloxTest,
	}

	valloxDumpCmd = &cobra.Command{
		Use:     "dump",
		Aliases: []string{"d"},
		Short:   "Dump incoming messages from the Vallox RS485 bus",
		Args:    cobra.NoArgs,
		Run:     runValloxDump,
	}
)

func runValloxTest(_ *cobra.Command, _ []string) {
	log.Info().Str("PORT", viper.GetString("vallox.port")).Msg("Waiting for serial connection")
	mode := &serial.Mode{
		BaudRate: viper.GetInt("vallox.baudRate"),
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	_, err := serial.Open(viper.GetString("vallox.port"), mode)
	if err != nil {
		log.Fatal().Err(err).Msg("Serial connection failed")
	}
	log.Info().Msg("Serial connection works!!!")
}

func runValloxDump(_ *cobra.Command, _ []string) {
	opts := vallox.NewClientOptions().
		SetSerialPort(viper.GetString("vallox.port")).
		SetDefaultMessageHandler(func(msg vallox.Message) {
			fmt.Printf("Got message %v\n", msg)
			// for i := 0; i < len(msg.Msg); i++ {
			// 	fmt.Printf("0x%x ", msg.Msg[i])
			// }
			// fmt.Println()
		})
	client := vallox.NewClient(opts)
	if err := client.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Serial connection failed")
	}
	defer client.Disconnect()

	client.StartListener()
}

func init() {
	valloxCmd.AddCommand(valloxTestCmd)
	valloxCmd.AddCommand(valloxDumpCmd)
	rootCmd.AddCommand(valloxCmd)
}
