package httpparser

import (
	"regexp"
	"strconv"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"

	methodRegex  = `^GET|POST|PUT|DELETE`
	uriRegex     = `/[a-z]+/[a-z]+/[0-9]+`
	clengthRegex = `Content-Length:\s([0-9]+)`
)

type HttpHeader struct {
	Method        string
	Uri           string
	ContentLength int
}

type HttpFrame struct {
	Header HttpHeader
	Body   []byte
}

func (hh HttpHeader) IsValid() bool {
	return len(hh.Method) > 0 && len(hh.Uri) > 0
}

func find(pattern string, s string) string {
	r := regexp.MustCompile(pattern)
	return r.FindString(s)
}

func findSubmatch(pattern string, s string) []string {
	r := regexp.MustCompile(pattern)
	return r.FindStringSubmatch(s)
}

func matchs(pattern string, s string) bool {
	r := regexp.MustCompile(pattern)
	return r.MatchString(s)
}

func GetMethod(s string) string {
	return find(`^GET|POST|PUT|DELETE`, s)
}

func GetURI(s string) string {
	return find(`/[a-z]+/[a-z]+/[0-9]+`, s)
}

func GetContentLength(s string) int {
	if !matchs(clengthRegex, s) {
		return 0
	}
	l, _ := strconv.Atoi(findSubmatch(clengthRegex, s)[1])
	return l
}

func GetHeader(s string) HttpHeader {
	return HttpHeader{GetMethod(s), GetURI(s), GetContentLength(s)}
}
