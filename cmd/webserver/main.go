package main

import (
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/httpserver"
	"github.com/manuporto/distributedHTTPServer/pkg/util"
	"net"
	"os"
	"strconv"

	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
)

func calculateDestServer(uri string, dbServerName string, port string, numOfDbServers int) string {
	uriHash := util.CalculateHash(uri)
	dbnum := strconv.Itoa(int(uriHash)%numOfDbServers + 1)
	return dbServerName + dbnum + port
}

func fowardRequest(req *httpparser.HttpFrame, address string) (*httpparser.HttpFrame, error) {
	c, err := net.Dial("tcp4", address)
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

	dbServerName := os.Getenv("DBSRVNAME")
	numOfDbServers, _ := strconv.Atoi(os.Getenv("DBSRVS"))
	destServer := calculateDestServer(
		req.Header.URI,
		dbServerName,
		os.Getenv("DBSRVPORT"),
		numOfDbServers)
	fmt.Println("Destination: ", destServer)
	res, err := fowardRequest(req, destServer)
	c.Write([]byte(res.Raw))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of args\n Usage: ./dbserver <address>")
		return
	}
	sv := server.NewServer(os.Args[1], handleConnection)
	sv.Serve()
}
