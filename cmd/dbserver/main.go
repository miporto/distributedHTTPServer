package main

import (
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Wrong number of args\n Usage: ./dbserver <address> <lock pool size> <cache size>")
		return
	}
	lockpoolSize, _ := strconv.Atoi(os.Args[2])
	cacheSize, _ := strconv.Atoi(os.Args[3])
	dbs := NewDbServer(uint(lockpoolSize), uint(cacheSize))
	sv := server.NewServer(os.Args[1], &dbs)
	sv.Serve()
}
