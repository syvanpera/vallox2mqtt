package cmd

import (
	"fmt"
	"html/template"
	"os"

	"github.com/denisbrodbeck/machineid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/vallox2mqtt/ha"
)

var (
	mqttCmd = &cobra.Command{
		Use:     "mqtt",
		Aliases: []string{"m"},
		Short:   "MQTT commands",
	}

	mqttTestCmd = &cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Test MQTT connection",
		Args:    cobra.NoArgs,
		Run:     runMQTTTest,
	}

	mqttDiscoverCmd = &cobra.Command{
		Use:     "discover",
		Aliases: []string{"disc", "d"},
		Short:   "Introduce device and sensors to HA via MQTT device discovery",
		Args:    cobra.NoArgs,
		Run:     runMQTTDiscover,
	}
)

func runMQTTTest(_ *cobra.Command, _ []string) {
	log.Info().Str("BROKER", viper.GetString("MQTT.broker")).Msg("Connecting to MQTT")
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s", viper.GetString("MQTT.broker"))).
		SetClientID("vallox2mqtt").
		SetUsername(viper.GetString("MQTT.username")).
		SetPassword(viper.GetString("MQTT.password"))

	opts.OnConnect = func(_ mqtt.Client) {
		log.Info().Msg("MQTT connection works!!!")
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("MQTT connection failed")
	}

	client.Disconnect(250)
}

func runMQTTDiscover(_ *cobra.Command, _ []string) {
	id, err := machineid.ID()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get machine id")
	}

	tpl := template.Must(template.ParseFiles("templates/mqtt_discovery.tmpl"))

	device := ha.Device{
		Version:    "0.0.1",
		DeviceId:   id,
		DeviceName: "vallox_digit_se",
		SensorId:   "fan_speed",
		SensorName: "Fan Speed",
		SensorIcon: "fan",
	}

	err = tpl.Execute(os.Stdout, device)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to execute template")
	}
}

func init() {
	mqttCmd.AddCommand(mqttTestCmd)
	mqttCmd.AddCommand(mqttDiscoverCmd)
	rootCmd.AddCommand(mqttCmd)
}
