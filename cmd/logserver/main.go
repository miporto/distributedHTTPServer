package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/manuporto/distributedHTTPServer/pkg/logger"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	var msgLen uint32
	for {
		err := binary.Read(c, binary.LittleEndian, msgLen)
		if err != nil {
			// send error msg
		}
		msg := make([]byte, msgLen)
		read, err := c.Read(msg)
		if read < len(msg) || err != nil {
			// send error msg
		}
		// chequear si es necesario mandar mensaje de ok
		fmt.Println(msg) //loggear
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Wrong number of args\n Usage: ./logserver <address> <log-file>")
		return
	}
	logger.NewLogger(os.Args[2])
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
