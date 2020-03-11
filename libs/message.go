package libs

type MessageType byte

const (
	Ping MessageType = '0'
	Pong MessageType = '1'

	Login MessageType = '3'
	LoginResp MessageType = '4'

	NewProxy MessageType = '5'
	NewProxyResp MessageType = '6'
	CloseProxy MessageType = '7'

	ReqWorkConn	MessageType = '8'
	NewWorkConn MessageType = '8'
	StartWorkConn MessageType = 'a'
)

type MessageInterface interface {
	Dumps() []byte
	Type() MessageType
	Len() int64
}

type message struct {
	Header MessageType
	Length int64
	Bytes  []byte
}

func (m message) Dumps() []byte {
	return nil
}

func (m message) Type() MessageType {
	return m.Header
}

func (m message) Len() int64 {
	return m.Length
}

func NewMessage(t MessageType, bytes []byte) MessageInterface {
	length := 0
	if bytes != nil {
		length = len(bytes)
	}
	return message{
		Header: t,
		Length: int64(length),
		Bytes:  bytes,
	}
}
