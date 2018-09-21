package server

import (
	"fmt"
	"net"
)

type Server struct {
	Address string
	Handler func(net.Conn)
}

func NewServer(address string, handler func(net.Conn)) Server {
	return Server{address, handler}
}

func (s Server) Serve() {

	l, err := net.Listen("tcp4", s.Address)
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
		go s.Handler(c)
	}
}
