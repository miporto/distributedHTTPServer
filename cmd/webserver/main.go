package main

import (
	"fmt"
	"log"
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
		return nil, err
	}
	defer c.Close()
	c.Write([]byte(req.Raw))
	return httpserver.ReadRequest(c)
}

func handleConnection(c net.Conn) {
	defer c.Close()
	log2 := logger.GetInstance()
	log2.Info(fmt.Sprintf("Serving %s", c.RemoteAddr().String()))
	log.Printf("Serving %s", c.RemoteAddr().String())
	req, err := httpserver.ReadRequest(c)
	if err != nil || !req.IsValid() {
		log.Println("Closing due to invalid HTTP request: ", req.Raw)
		log2.Error(fmt.Sprintf("Closing due to invalid HTTP request: %s from %s",
			req.Raw, c.RemoteAddr().String()))
		httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: httpparser.StatusBadRequest})
		return
	}

	dbServerName := os.Getenv("DBSRVNAME")
	numOfDbServers, _ := strconv.Atoi(os.Getenv("DBSRVS"))
	destServer := calculateDestServer(
		req.Header.URI,
		dbServerName,
		os.Getenv("DBSRVPORT"),
		numOfDbServers)
	log2.Info(fmt.Sprintf("HTTP Request: %s with destination: %s", req.Raw, destServer))
	log.Printf("HTTP Request: %s with destination: %s", req.Raw, destServer)
	res, err := fowardRequest(req, destServer)
	if err != nil {
		log2.Error(err.Error())
		return
	}
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
