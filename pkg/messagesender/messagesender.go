package messagesender

import (
	"encoding/binary"
	"fmt"
	"net"
)

type MessageSender struct {
	address string
	msgCh   <-chan string
}

func NewMessageSender(address string, msgCh <-chan string) *MessageSender {
	return &MessageSender{address: address, msgCh: msgCh}
}

func (ms *MessageSender) start() {
	c, err := net.Dial("tcp4", ms.address)
	if err != nil {
		fmt.Println(err)
		return
	}
	for msg := range ms.msgCh {
		binary.Write(c, binary.LittleEndian, len(msg)) // check error
		c.Write([]byte(msg))                           //error
	}
}
