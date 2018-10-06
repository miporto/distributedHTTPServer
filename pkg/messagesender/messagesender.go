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

func (ms *MessageSender) Start() {
	c, err := net.Dial("tcp4", "distributedhttpserver_log_1:8082")
	if err != nil {
		fmt.Println("MessageSender error: ", err)
		return
	}
	defer c.Close()
	for msg := range ms.msgCh {
		msgLen := uint32(len(msg))
		fmt.Println("Sending ", msg, " with length ", msgLen)
		err := binary.Write(c, binary.LittleEndian, msgLen) // check error
		if err != nil {
			fmt.Println("Error when sending size: ", err)
		}
		n, err := c.Write([]byte(msg)) //error
		if n < int(msgLen) || err != nil {
			fmt.Println("Error when sending message: ", err, "bytes sent: ", n)
		}
	}
}
