package httpparser

import (
	"regexp"
	"strconv"
)

type HttpHeader struct {
	method        string
	uri           string
	contentLength int
}

type HttpFrame struct {
	header HttpHeader
	body   []byte
}

func find(pattern string, s string) string {
	r := regexp.MustCompile(pattern)
	return r.FindString(s)
}

func findSubmatch(pattern string, s string) []string {
	r := regexp.MustCompile(pattern)
	return r.FindStringSubmatch(s)
}
func GetMethod(s string) string {
	return find(`^[A-Z]+\b`, s)
}

func GetURI(s string) string {
	return find(`/[a-z]+/[a-z]+/[0-9]+`, s)
}

func GetHost(s string) string {
	return find(`^Host:\s([a-z]+:*[0-9]*)`, s)
}

func GetContentLength(s string) int {
	l, _ := strconv.Atoi(findSubmatch(`Content-Length:\s([0-9]+)`, s)[1])
	return l
}

func GetHeader(s string) HttpHeader {
	return HttpHeader{GetMethod(s), GetURI(s), GetContentLength(s)}
}
