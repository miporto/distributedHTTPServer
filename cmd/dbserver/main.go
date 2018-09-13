package main

import (
	"fmt"
	"github.com/manuporto/distributedHTTPServer/pkg/filemanager"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		title := r.URL.Path[1:] + ".json"
		b, _ := filemanager.LoadFile(title)
		fmt.Fprintf(w, "<h1>File</h1><div>%s</div>", b)
	case http.MethodDelete:
		filemanager.DeleteFile(r.URL.Path[1:])

	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
