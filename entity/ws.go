package entity

type WSMessageType int

const (
	WSMessageTypeCreateMessageSuccess WSMessageType = iota + 1
	WSMessageTypeNewMessage
)

type WSMessage struct {
	Type    WSMessageType
	Payload interface{}
}
