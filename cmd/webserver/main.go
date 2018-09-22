package main

import (
	"fmt"
	"net"

	"github.com/manuporto/distributedHTTPServer/pkg/ftpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/httpserver"

	"github.com/manuporto/distributedHTTPServer/pkg/ftpclient"
	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
)

func fowardRequest(req *httpparser.HttpFrame) *ftpparser.FTPResponse {
	ftpClient, _ := ftpclient.Connect("address")
	defer ftpClient.Close()
	ftpClient.Send(ftpparser.HTTPToFTP(req))
	return ftpClient.Receive()
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
	fowardRequest(req)
	c.Write([]byte("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n"))
}

func main() {
	sv := server.NewServer(":8080", handleConnection)
	sv.Serve()
}
