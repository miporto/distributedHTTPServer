package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/manuporto/distributedHTTPServer/pkg/filemanager"
	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"

	"github.com/manuporto/distributedHTTPServer/pkg/httpserver"
)

func handleGET(c net.Conn, req *httpparser.HttpFrame, fm *filemanager.FileManager) {
	status := httpparser.StatusOK
	body, err := fm.Load(req.Header.URI)
	if err != nil {
		if os.IsNotExist(err) {
			status = httpparser.StatusNotFound
		} else {
			status = httpparser.StatusInternalServerError
			body = []byte(err.Error())
		}
		fmt.Println("Error in get: ", err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{
		Status:        status,
		ContentType:   httpparser.JSONContentType,
		ContentLength: len(body),
		Body:          body})

}

func handlePOST(c net.Conn, req *httpparser.HttpFrame, fm *filemanager.FileManager) {
	status := httpparser.StatusOK
	var body []byte
	err := fm.Save(req.Header.URI, req.Body)
	if err != nil {
		if os.IsExist(err) {
			status = httpparser.StatusConflict
		} else {
			status = httpparser.StatusInternalServerError
		}
		body = []byte(err.Error())
		fmt.Println("Error in post: ", err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: status, Body: body})
}

func handlePUT(c net.Conn, req *httpparser.HttpFrame, fm *filemanager.FileManager) {
	status := httpparser.StatusOK
	var body []byte
	err := fm.Update(req.Header.URI, req.Body)
	if err != nil {
		if os.IsNotExist(err) {
			status = httpparser.StatusNotFound
		} else {
			status = httpparser.StatusInternalServerError
			body = []byte(err.Error())
		}
		fmt.Println("Error in put: ", err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: status, Body: body})
}

func handleDELETE(c net.Conn, req *httpparser.HttpFrame, fm *filemanager.FileManager) {
	status := httpparser.StatusOK
	var body []byte
	err := fm.Delete(req.Header.URI)
	if err != nil {
		if os.IsNotExist(err) {
			status = httpparser.StatusNotFound
			fmt.Println(status)
		} else {
			status = httpparser.StatusInternalServerError
			body = []byte(err.Error())
		}
		fmt.Println("Error in delete: ", err)
	}
	httpserver.WriteResponse(c, &httpparser.HttpResponse{Status: status, Body: body})
}

func handleConnection(c net.Conn, fm *filemanager.FileManager) {
	defer c.Close()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	req, err := httpserver.ReadRequest(c)
	if err != nil {
		// TODO handle error
		fmt.Println("Closing: ", err)
		c.Write([]byte(err.Error()))
		return
	}
	req.Header.URI = req.Header.URI[1:] + ".json"

	switch req.Header.Method {
	case httpparser.MethodGet:
		handleGET(c, req, fm)
	case httpparser.MethodPost:
		handlePOST(c, req, fm)
	case httpparser.MethodPut:
		handlePUT(c, req, fm)
	case httpparser.MethodDelete:
		handleDELETE(c, req, fm)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Wrong number of args\n Usage: ./dbserver <address>")
		return
	}
	lockpoolSize, _ := strconv.Atoi(os.Args[2])
	cacheSize, _ := strconv.Atoi(os.Args[3])
	l, err := net.Listen("tcp4", os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fm := filemanager.NewFileManager(uint(lockpoolSize), uint(cacheSize))
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c, &fm)
	}
}
