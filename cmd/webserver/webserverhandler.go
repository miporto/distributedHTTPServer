package main

import (
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/httpserver"
	"github.com/manuporto/distributedHTTPServer/pkg/util"
	"log"
	"net"
	"os"
	"strconv"
)

type WebServer struct {
	logCh chan<- string
}

func NewWebServer(logCh chan<- string) WebServer {
	return WebServer{logCh: logCh}
}

func (ws *WebServer) Handle(c net.Conn) {
	defer c.Close()
	log.Printf("Serving %s", c.RemoteAddr().String())
	ws.logCh <- fmt.Sprintf("Serving %v", c.RemoteAddr())
	req, err := httpserver.ReadRequest(c)
	if err != nil || !req.IsValid() {
		log.Println("Closing due to invalid HTTP request: ", req.Raw)
		ws.logCh <- fmt.Sprintf("Closing due to invalid HTTP request: %s from %v", req.Raw, c.RemoteAddr())
		httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: httpparser.StatusBadRequest})
		return
	}

	dbServerName := os.Getenv("DBSRVNAME")
	numOfDbServers, _ := strconv.Atoi(os.Getenv("DBSRVS"))
	destServer := ws.calculateDestServer(
		req.Header.URI,
		dbServerName,
		os.Getenv("DBSRVPORT"),
		numOfDbServers)
	ws.logCh <- fmt.Sprintf("HTTP Request: %s with destination: %s", req.Raw, destServer)
	log.Printf("HTTP Request: %s with destination: %s", req.Raw, destServer)
	res, err := ws.fowardRequest(req, destServer)
	if err != nil {
		ws.logCh <- err.Error()
		return
	}
	c.Write([]byte(res.Raw))
}

func (ws *WebServer) calculateDestServer(uri string, dbServerName string, port string, numOfDbServers int) string {
	if dbServerName == "localhost" {
		return dbServerName + port
	}
	uriHash := util.CalculateHash(uri)
	dbnum := strconv.Itoa(int(uriHash)%numOfDbServers + 1)
	return dbServerName + dbnum + port
}

func (ws *WebServer) fowardRequest(req *httpparser.HttpFrame, address string) (*httpparser.HttpFrame, error) {
	c, err := net.Dial("tcp4", address)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	c.Write([]byte(req.Raw))
	return httpserver.ReadRequest(c)
}
