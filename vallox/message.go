package vallox

type MessageHandler func(Message)

type Message struct {
	Msg []byte
}

func NewMessage(msg []byte) Message {
	return Message{
		Msg: msg,
	}
}
