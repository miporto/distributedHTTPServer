package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	address string
	handler Handler
}

type Handler interface {
	Handle(net.Conn)
}

func NewServer(address string, handler Handler) Server {
	return Server{address, handler}
}

func (s Server) Serve() {
	l, err := net.Listen("tcp4", s.address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	wg := &sync.WaitGroup{}
	defer wg.Wait()
	maxWorkers := 100
	guard := make(chan struct{}, maxWorkers)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.handler.Handle(c)
			<-guard
		}()
	}
}
