package main

import (
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/filemanager"
	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/httpserver"
	"log"
	"net"
	"os"
)

type DbServer struct {
	fm *filemanager.FileManager
}

func NewDbServer(lockpoolSize uint, cacheSize uint) DbServer {
	fm := filemanager.NewFileManager(lockpoolSize, cacheSize)
	return DbServer{fm: &fm}
}

func (dbs *DbServer) handleGET(c net.Conn, req *httpparser.HttpFrame) {
	status := httpparser.StatusOK
	body, err := dbs.fm.Load(req.Header.URI)
	if err != nil {
		if os.IsNotExist(err) {
			status = httpparser.StatusNotFound
		} else {
			status = httpparser.StatusInternalServerError
			body = []byte(err.Error())
		}
		log.Printf("ERROR: [Conn %s] %s\n", c.RemoteAddr().String(), err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{
		Status:        status,
		ContentType:   httpparser.JSONContentType,
		ContentLength: len(body),
		Body:          body})

}

func (dbs *DbServer) handlePOST(c net.Conn, req *httpparser.HttpFrame) {
	status := httpparser.StatusOK
	var body []byte
	err := dbs.fm.Save(req.Header.URI, req.Body)
	if err != nil {
		if os.IsExist(err) {
			status = httpparser.StatusConflict
		} else {
			status = httpparser.StatusInternalServerError
		}
		body = []byte(err.Error())
		log.Printf("ERROR: [Conn %s] %s\n", c.RemoteAddr().String(), err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: status, Body: body})
}

func (dbs *DbServer) handlePUT(c net.Conn, req *httpparser.HttpFrame) {
	status := httpparser.StatusOK
	var body []byte
	err := dbs.fm.Update(req.Header.URI, req.Body)
	if err != nil {
		if os.IsNotExist(err) {
			status = httpparser.StatusNotFound
		} else {
			status = httpparser.StatusInternalServerError
			body = []byte(err.Error())
		}
		log.Printf("ERROR: [Conn %s] %s\n", c.RemoteAddr().String(), err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: status, Body: body})
}

func (dbs *DbServer) handleDELETE(c net.Conn, req *httpparser.HttpFrame) {
	status := httpparser.StatusOK
	var body []byte
	err := dbs.fm.Delete(req.Header.URI)
	if err != nil {
		if os.IsNotExist(err) {
			status = httpparser.StatusNotFound
			fmt.Println(status)
		} else {
			status = httpparser.StatusInternalServerError
			body = []byte(err.Error())
		}
		log.Printf("ERROR: [Conn %s] %s\n", c.RemoteAddr().String(), err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: status, Body: body})
}

func (dbs *DbServer) Handle(c net.Conn) {
	defer c.Close()
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	req, err := httpserver.ReadRequest(c)
	if err != nil {
		log.Println("Closing due to unexpected error: ", err)
		httpserver.WriteResponse(c, &httpparser.HttpResponse{
			Status: httpparser.StatusInternalServerError,
			Body:   []byte(err.Error())})
		return
	}
	log.Println(req.Raw)
	if !req.IsValid() {
		log.Println("Closing due to invalid HTTP request: ", req.Header.URI)
		httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: httpparser.StatusBadRequest})
		return
	}
	req.Header.URI = req.Header.URI + ".json"

	switch req.Header.Method {
	case httpparser.MethodGet:
		dbs.handleGET(c, req)
	case httpparser.MethodPost:
		dbs.handlePOST(c, req)
	case httpparser.MethodPut:
		dbs.handlePUT(c, req)
	case httpparser.MethodDelete:
		dbs.handleDELETE(c, req)
	}
}
