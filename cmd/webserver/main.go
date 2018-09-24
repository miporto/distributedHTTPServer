package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/manuporto/distributedHTTPServer/pkg/httpserver"
	"github.com/manuporto/distributedHTTPServer/pkg/logger"
	"github.com/manuporto/distributedHTTPServer/pkg/util"

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
	log := logger.GetInstance()
	log.Info(fmt.Sprintf("Serving %s", c.RemoteAddr().String()))
	req, err := httpserver.ReadRequest(c)
	if err != nil {
		// TODO handle error
		log.Error(fmt.Sprintf("Closing %s because of error: %s", c.RemoteAddr().String(), err.Error()))
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
	log.Info(fmt.Sprintf("HTTP Request: %s withh destination: %s", req.Raw, destServer))
	res, err := fowardRequest(req, destServer)
	c.Write([]byte(res.Raw))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of args\n Usage: ./dbserver <address>")
		return
	}
	logger.NewLogger("log.txt")
	sv := server.NewServer(os.Args[1], handleConnection)
	sv.Serve()
}
