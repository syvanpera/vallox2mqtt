package cmd

import (
	"bytes"
	"fmt"
	"html/template"

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
		Run:     runTest,
	}

	mqttDiscoverCmd = &cobra.Command{
		Use:     "discover",
		Aliases: []string{"disc", "d"},
		Short:   "Introduce device and sensors to HA via MQTT device discovery",
		Args:    cobra.NoArgs,
		Run:     runDiscovery,
	}
)

func runTest(_ *cobra.Command, _ []string) {
	client := mqttConnect()
	defer client.Disconnect(250)
}

func runDiscovery(_ *cobra.Command, _ []string) {
	client := mqttConnect()
	defer client.Disconnect(250)

	sendDiscovery(client, "sensor", "fan_speed", "Fan Speed", "", "")
	sendDiscovery(client, "sensor", "inside_temp", "Inside Temperature", "temperature", "°C")
	sendDiscovery(client, "sensor", "outside_temp", "Outside Temperature", "temperature", "°C")
}

func mqttConnect() mqtt.Client {
	log.Info().Str("BROKER", viper.GetString("MQTT.broker")).Msg("Connecting to MQTT")
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s", viper.GetString("MQTT.broker"))).
		SetClientID("vallox2mqtt").
		SetUsername(viper.GetString("MQTT.username")).
		SetPassword(viper.GetString("MQTT.password"))

	opts.OnConnect = func(c mqtt.Client) {
		log.Info().Msg("MQTT broker connected")
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("MQTT connection failed")
	}

	return client
}

func sendDiscovery(client mqtt.Client, component, entity_id, entity_name, device_class, unit_of_measurement string) {
	tpl := template.Must(template.ParseFiles("templates/mqtt_discovery.tmpl"))

	device_id := viper.GetString("ha.device_id")

	device := ha.Device{
		Version:           "0.0.1",
		DeviceId:          device_id,
		DeviceName:        viper.GetString("ha.device_name"),
		DeviceClass:       device_class,
		EntityId:          entity_id,
		EntityName:        entity_name,
		UnitOfMeasurement: unit_of_measurement,
	}

	var discovery bytes.Buffer
	if err := tpl.Execute(&discovery, device); err != nil {
		log.Fatal().Err(err).Msg("Unable to execute template")
	}

	prefix := viper.GetString("ha.discovery_prefix")
	topic := fmt.Sprintf("%s/%s/%s/%s/config", prefix, component, device_id, entity_id)
	log.Debug().Str("TOPIC", topic).Msg("MQTT Publish")
	token := client.Publish(topic, 0, false, discovery)
	token.Wait()
}

func init() {
	mqttCmd.AddCommand(mqttTestCmd)
	mqttCmd.AddCommand(mqttDiscoverCmd)
	rootCmd.AddCommand(mqttCmd)
}
