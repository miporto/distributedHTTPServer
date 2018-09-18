package httpserver

import (
	"github.com/manuporto/distributedHTTPServer/pkg/server"
	"net"
)

type HttpServer struct {
	s server.Server
}

func (hs HttpServer) ListenAndServe(address string, handler func(net.Conn)) {
	hs.s = server.Server{address, handler}
	hs.s.Serve()
}
