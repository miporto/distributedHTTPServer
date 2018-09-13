package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type File struct {
	Name string
	Body []byte
}

func (f *File) save() error {
	filename := f.Name + ".txt"
	return ioutil.WriteFile(filename, f.Body, 0600)
}

func loadFile(title string) ([]byte, error) {
	filename := title + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	b, _ := loadFile(title)
	fmt.Fprintf(w, "<h1>File</h1><div>%s</div>", b)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
