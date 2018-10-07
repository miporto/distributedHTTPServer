package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	origin = "origin_"
	entity = "entity_"
)

func generateRandomURI() string {
	var uri strings.Builder
	originNum := strconv.Itoa(rand.Intn(10))
	entityNum := strconv.Itoa(rand.Intn(5))
	id := strconv.Itoa(rand.Intn(10))
	uri.WriteString("/")
	uri.WriteString(origin)
	uri.WriteString(originNum)
	uri.WriteString("/")
	uri.WriteString(entity)
	uri.WriteString(entityNum)
	uri.WriteString("/")
	uri.WriteString(id)

	return uri.String()
}

func getWorker(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		start := time.Now()
		fullUrl := url + generateRandomURI()
		resp, err := http.Get(fullUrl)
		if err != nil {
			log.Printf("ERROR: [GET %s] failed due to %s", fullUrl, err.Error())
			return
		}
		defer resp.Body.Close()

		secs := time.Since(start).Seconds()
		log.Printf("INFO: [GET %s] Elapsed: %v response status: %s", fullUrl, secs, resp.Status)
	}
}

func postWorker(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		start := time.Now()
		fullUrl := url + generateRandomURI()
		var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR: [POST %s] failed due to %s", fullUrl, err.Error())
			return
		}
		defer resp.Body.Close()

		secs := time.Since(start).Seconds()
		log.Printf("INFO: [POST %s] Elapsed: %v response status: %s", fullUrl, secs, resp.Status)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of args\n Usage: ./loadtester <address>")
		return
	}

	var wg sync.WaitGroup
	workers := 200
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go postWorker(os.Args[1], &wg)
	}
	wg.Wait()

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go getWorker(os.Args[1], &wg)
	}
	wg.Wait()
}
