package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "db:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection established, sending packets")
	fmt.Fprintf(conn, "Hello\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
	fmt.Fprintf(conn, "STOP\n")
	status, err = bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
}
