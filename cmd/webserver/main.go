package main

import (
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/messagesender"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of args\n Usage: ./dbserver <address>")
		return
	}
	logCh := make(chan string, 1000)
	defer close(logCh)
	ms := messagesender.NewMessageSender(os.Getenv("LOGSRV"), logCh)
	go ms.Start()
	ws := NewWebServer(logCh)
	sv := server.NewServer(os.Args[1], &ws)
	sv.Serve()
}
