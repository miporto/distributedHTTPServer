package server

import (
	"fmt"
	"net"
)

type Server struct {
	address string
	handler func(net.Conn)
}

func (s Server) Serve() {

	l, err := net.Listen("tcp4", s.address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go s.handler(c)
	}
}
