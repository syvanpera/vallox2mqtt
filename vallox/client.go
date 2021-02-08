package vallox

import (
	"errors"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
	"go.bug.st/serial"
)

type Client interface {
	Connect() error
	IsConnected() bool
	Disconnect()
}

type client struct {
	port       serial.Port
	options    ClientOptions
	serialMode serial.Mode
	connected  bool
	backOff    backoff.BackOff
}

var ErrConnectionFailed = errors.New("Could not open serial connection")
var ErrConnectionLost = errors.New("Serial connection lost")
var ErrNotConnected = errors.New("Not Connected")

func NewClient(o *ClientOptions) Client {
	c := &client{}
	c.options = *o
	c.serialMode = serial.Mode{
		BaudRate: c.options.BaudRate,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	c.backOff = c.backOffStrategy()

	return c
}

func (c *client) Connect() error {
	if c.options.WaitForConnection {
		log.Info().Str("PORT", c.options.SerialPort).Msg("Waiting for serial connection")
	} else {
		log.Info().Str("PORT", c.options.SerialPort).Msg("Opening serial connection")
	}

	if err := backoff.Retry(c.openSerialConnection, c.backOff); err != nil {
		log.Error().Err(err).Msg("Serial connection failed")
		return err
	}

	c.connected = true

	if err := c.startListener(); err != nil {
		log.Error().Err(err).Msg("Serial connection lost")
		c.Disconnect()
		if c.options.AutoReconnect {
			log.Info().Msg("Reconnecting")
			c.Connect()
		}
	}

	return nil
}

func (c *client) IsConnected() bool {
	return c.connected
}

func (c *client) Disconnect() {
	c.connected = false
	if err := c.port.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close the serial connection")
	}
}

func (c *client) backOffStrategy() backoff.BackOff {
	if c.options.WaitForConnection {
		b := backoff.NewExponentialBackOff()
		b.MaxInterval = 30 * time.Second
		b.Reset()

		return b
	} else {
		b := backoff.StopBackOff{}
		b.Reset()

		return &b
	}
}

func (c *client) openSerialConnection() error {
	var err error
	c.port, err = serial.Open(c.options.SerialPort, &c.serialMode)

	return err
}

func (c *client) startListener() error {
	log.Info().Msg("Starting Vallox message listener")

	buff := make([]byte, valloxMessageLength)
	for {
		n, err := c.port.Read(buff)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while reading from serial connection")
		}
		switch {
		case n == 0:
			return ErrConnectionLost
		case n < valloxMessageLength:
			log.Warn().Msg("Incomplete or invalid message received")
			continue
		}

		domain := buff[0]
		sender := buff[1]
		receiver := buff[2]
		command := buff[3]
		arg := buff[4]
		checksum := buff[5]
		computedChecksum := (domain + sender + receiver + command + arg) & 0x00ff

		if domain != valloxDomain {
			log.Warn().Int("DOMAIN", int(buff[0])).Msg("Unknown message")
			continue
		}

		if checksum == computedChecksum {
			c.parseMessage(sender, receiver, command, arg)
			c.options.DefaultMessageHandler(NewMessage(buff[1:n]))
		}
	}
}

func (c *client) parseMessage(sender byte, receiver byte, command byte, arg byte) {
	log.Debug().Int("FANSPEED", int(VariableFanSpeed)).Int("SENDER", int(sender)).Int("RECEIVER", int(receiver)).Int("COMMAND", int(command)).Int("ARG", int(arg)).Msg("Got message")

	if receiver == AddressPanels || receiver == AddressPanel1 {
		switch int(command) {
		case VariableFanSpeed:
			fanSpeed := convertFanSpeed(arg)
			log.Debug().Int("FANSPEED", fanSpeed).Msg("Got fan speed")
		case VariableTempOutside:
			temp := convertTemperature(arg)
			log.Debug().Int("TEMPERATURE", temp).Msg("Got outside air temperature")
		case VariableTempExhaust:
			temp := convertTemperature(arg)
			log.Debug().Int("TEMPERATURE", temp).Msg("Got exhaust air temperature")
		case VariableTempInside:
			temp := convertTemperature(arg)
			log.Debug().Int("TEMPERATURE", temp).Msg("Got inside air temperature")
		case VariableTempSupply:
			temp := convertTemperature(arg)
			log.Debug().Int("TEMPERATURE", temp).Msg("Got supply air temperature")
		}
	}
}
