package main

import (
	"encoding/binary"
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/logger"
	"log"
	"net"
	"os"
)

const maxsize = 2048

func handleConnection(c net.Conn) {
	defer c.Close()
	var msgLen uint32
	logger := logger.GetInstance()
	for {
		err := binary.Read(c, binary.LittleEndian, &msgLen)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in first receive: %s", err))
			// send error msg
		}
		fmt.Println("Received: ", msgLen)
		var msg []byte
		if msgLen <= maxsize {
			msg = make([]byte, msgLen)
		} else {
			msg = make([]byte, maxsize)
		}
		_, err = c.Read(msg)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in second receive: %s", err))
		}
		fmt.Println("Message: ", string(msg))
		logger.Info(string(msg))
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Wrong number of args\n Usage: ./logserver <address>")
	}
	l, err := net.Listen("tcp4", os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	handleConnection(c)
}
