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

func readHeader(c net.Conn) (string, error) {
	r := bufio.NewReader(c)
	b, err := r.ReadString('\n')
	var header strings.Builder
	header.WriteString(b)
	for {
		fmt.Print(b)
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

func readBody(c net.Conn, hh httpparser.HttpHeader) ([]byte, error) {
	fmt.Println("Making body of ", hh.ContentLength)
	body := make([]byte, hh.ContentLength)
	_, err := bufio.NewReader(c).Read(body)
	return body, err
}

func readRequest(c net.Conn) (*httpparser.HttpFrame, error) {
	header, err := readHeader(c)
	fmt.Print(header)
	if err != nil {
		return nil, err
	}
	httpheader := httpparser.GetHeader(header)
	fmt.Println("reading body...")
	body, err := readBody(c, httpheader)
	fmt.Println("Body read!")
	return &httpparser.HttpFrame{httpheader, body}, err
}

func handleGET(hf *httpparser.HttpFrame) []byte {
	fmt.Println(hf.Header.Uri)
	f, err := filemanager.LoadFile(hf.Header.Uri[1:])
	if err != nil {
		fmt.Println("Invalid file: ", err)
		return nil
	}
	return f
}

func handlePOST(hf *httpparser.HttpFrame) error {
	fmt.Println("Handling POST...")
	return filemanager.SaveFile(hf.Header.Uri[1:], hf.Body)
}

func handlePUT() {

}

func handleDELETE() {

}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		req, err := readRequest(c)
		if err != nil {
			// TODO handle error
			c.Write([]byte(err.Error()))
		}
		fmt.Println(req.Header.Method)
		switch req.Header.Method {
		case httpparser.MethodGet:
			f := handleGET(req)
			fmt.Printf(string(f))
			c.Write(f)
		case httpparser.MethodPost:
			err = handlePOST(req)
			c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

		}
	}
	c.Close()
}

func main() {
	sv := server.Server{":8080", handleConnection}
	sv.Serve()
}
