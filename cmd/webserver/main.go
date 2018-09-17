package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
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

func handleConnection(c net.Conn) {
	for {
		header, err := readHeader(c)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if header == "STOP\n" {
			fmt.Println("Stopping...")
			break
		}
		fmt.Println(header)
		fmt.Println(httpparser.GetMethod(header))
		fmt.Println(httpparser.GetURI(header))
		fmt.Println(httpparser.GetContentLength(header))
	}
	c.Close()
}

func main() {
	l, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
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
