package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/manuporto/distributedHTTPServer/pkg/filemanager"
	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
)

func readHeader(r *bufio.Reader) (string, error) {
	b, err := r.ReadString('\n')
	var header strings.Builder
	header.WriteString(b)
	for {
		if err != nil {
			return "", err
		}

		if b == "\r\n" {
			return header.String(), nil
		}
		b, err = r.ReadString('\n')
		header.WriteString(b)
	}
}

func readBody(r *bufio.Reader, size int) ([]byte, error) {
	body := make([]byte, size)
	// TODO handle read < body
	_, err := r.Read(body)
	return body, err
}

func readRequest(c net.Conn) (*httpparser.HttpFrame, error) {
	r := bufio.NewReader(c)
	header, err := readHeader(r)
	fmt.Print(header)
	if err != nil {
		return nil, err
	}
	httpheader := httpparser.GetHeader(header)
	body, err := readBody(r, httpheader.ContentLength)
	return &httpparser.HttpFrame{Header: httpheader, Body: body}, err
}

func handleGET(hf *httpparser.HttpFrame) []byte {
	f, err := filemanager.LoadFile(hf.Header.Uri)
	if err != nil {
		fmt.Println("Invalid file: ", err)
		return nil
	}
	return f
}

func handlePOST(hf *httpparser.HttpFrame) error {
	return filemanager.SaveFile(hf.Header.Uri, hf.Body)
}

func handlePUT(hf *httpparser.HttpFrame) {
	filemanager.UpdateFile(hf.Header.Uri, hf.Body)
}

func handleDELETE(hf *httpparser.HttpFrame) {
	filemanager.DeleteFile(hf.Header.Uri)
}

func handleConnection(c net.Conn) {
	defer c.Close()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	req, err := readRequest(c)
	if err != nil {
		// TODO handle error
		fmt.Println("Closing: ", err)
		c.Write([]byte(err.Error()))
		return
	}
	req.Header.Uri = req.Header.Uri[1:] + ".json"
	switch req.Header.Method {
	case httpparser.MethodGet:
		f := handleGET(req)
		fmt.Printf(string(f))
	case httpparser.MethodPost:
		handlePOST(req)
	case httpparser.MethodPut:
		handlePUT(req)
	case httpparser.MethodDelete:
		handleDELETE(req)

	}
	c.Write([]byte("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n"))
}

func main() {
	sv := server.Server{Address: ":8080", Handler: handleConnection}
	sv.Serve()
}
