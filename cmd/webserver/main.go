package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "distributedhttpserver_db_1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection established, sending packets")
	time.Sleep(5 * time.Second)
	fmt.Fprintf(conn, "hello\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
	fmt.Fprintf(conn, "STOP\n")
	status, err = bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
}
