package main

import (
	"bufio"
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
	"net"
	"strings"
)

func readHeader(c net.Conn) (string, error) {
	r := bufio.NewReader(c)
	b, err := r.ReadString('\n')
	var header strings.Builder
	header.WriteString(b)
	for {
		if err != nil {
			return "", err
		}

		if b == "\n" {
			return header.String(), nil
		}
		b, err = r.ReadString('\n')
		header.WriteString(b)
	}
}

func handleGET(hf httpparser.HttpFrame) {

}

func handlePOST() {

}

func handlePUT() {

}

func handleDELETE() {

}

func handleConnection(c net.Conn) {
	for {
		header, err := readHeader(c)
		if err != nil {
			fmt.Println(err)
			break
		}
		if header == "STOP\n" {
			fmt.Println("Stopping...")
			break
		}
		fmt.Println(header)
		httpheader := httpparser.GetHeader(header)
		body := make([]byte, httpheader.ContentLength)
		read, err := bufio.NewReader(c).Read(body)

		if read < len(body) || err != nil {
			// TODO handle invalid http frame
			fmt.Println(err)
			continue
		}
		frame := httpparser.HttpFrame{httpheader, body}

		switch httpheader.Method {
		case httpparser.MethodGet:
			handleGET(frame)
		}
	}
	c.Close()
}

func main() {
	sv := server.Server{":8080", handleConnection}
	sv.Serve()
}

// func main() {
// 	conn, err := net.Dial("tcp", "distributedhttpserver_db_1:8080")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()
// 	fmt.Println("Connection established, sending packets")
// 	time.Sleep(5 * time.Second)
// 	fmt.Fprintf(conn, "hello\n")
// 	status, err := bufio.NewReader(conn).ReadString('\n')
// 	fmt.Println(status)
// 	fmt.Fprintf(conn, "STOP\n")
// 	status, err = bufio.NewReader(conn).ReadString('\n')
// 	fmt.Println(status)
// }
