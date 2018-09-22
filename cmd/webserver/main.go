package main

import (
	"fmt"
	"net"

	"github.com/manuporto/distributedHTTPServer/pkg/httpserver"

	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
)

func fowardRequest(req *httpparser.HttpFrame) (*httpparser.HttpFrame, error) {
	c, err := net.Dial("tcp4", ":8081")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer c.Close()
	c.Write([]byte(req.Raw))
	return httpserver.ReadRequest(c)
}

func handleConnection(c net.Conn) {
	defer c.Close()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	req, err := httpserver.ReadRequest(c)
	if err != nil {
		// TODO handle error
		fmt.Println("Closing: ", err)
		c.Write([]byte(err.Error()))
		return
	}
	res, err := fowardRequest(req)
	c.Write([]byte(res.Raw))
}

func main() {
	sv := server.NewServer(":8080", handleConnection)
	sv.Serve()
}
