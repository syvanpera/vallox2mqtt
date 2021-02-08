package vallox

type ClientOptions struct {
	SerialPort            string
	BaudRate              int
	WaitForConnection     bool
	AutoReconnect         bool
	DefaultMessageHandler MessageHandler
}

func NewClientOptions() *ClientOptions {
	o := &ClientOptions{
		SerialPort:        "/dev/ttyUSB0",
		BaudRate:          9600,
		WaitForConnection: true,
		AutoReconnect:     true,
	}
	return o
}

func (o *ClientOptions) SetSerialPort(port string) *ClientOptions {
	o.SerialPort = port
	return o
}

func (o *ClientOptions) SetBaudRate(rate int) *ClientOptions {
	o.BaudRate = rate
	return o
}

func (o *ClientOptions) SetWaitForConnection(w bool) *ClientOptions {
	o.WaitForConnection = w
	return o
}

func (o *ClientOptions) SetAutoReconnect(a bool) *ClientOptions {
	o.AutoReconnect = a
	return o
}

func (o *ClientOptions) SetDefaultMessageHandler(m MessageHandler) *ClientOptions {
	o.DefaultMessageHandler = m
	return o
}
