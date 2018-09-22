package httpparser

import (
	"regexp"
	"strconv"
)

const (
	Version = "HTTP/1.1"
	Delim   = "\r\n"

	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"

	StatusOK                  = "200 OK"
	StatusBadRequest          = "400 Bad Request"
	StatusNotFound            = "404 Not Found"
	StatusInternalServerError = "500 Internal Server Error"

	ContentTypeResHeader   = "Content-Type: "
	ContentLengthResHeader = "Content-Length: "
	ConnectionHeader       = "Connection: Close"

	JSONContentType = "application/json"

	methodRegex  = `^GET|POST|PUT|DELETE`
	uriRegex     = `/[a-z]+/[a-z]+/[0-9]+`
	clengthRegex = `Content-Length:\s([0-9]+)`
)

type HttpHeader struct {
	Status        string
	Method        string
	URI           string
	ContentType   string
	ContentLength int
}

type HttpFrame struct {
	Raw    string
	Header HttpHeader
	Body   []byte
}

type HttpResponse struct {
	Status        string
	ContentType   string
	ContentLength int
	Body          []byte
}

func (hh HttpHeader) IsValid() bool {
	return len(hh.Method) > 0 && len(hh.URI) > 0
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

func getStatus(s string) string {
	return ""
}

func GetMethod(s string) string {
	return find(methodRegex, s)
}

func GetURI(s string) string {
	return find(uriRegex, s)
}

func GetContentLength(s string) int {
	if !matchs(clengthRegex, s) {
		return 0
	}
	l, _ := strconv.Atoi(findSubmatch(clengthRegex, s)[1])
	return l
}

func GetHeader(s string) HttpHeader {
	return HttpHeader{
		Status:        "",
		Method:        GetMethod(s),
		URI:           GetURI(s),
		ContentType:   "",
		ContentLength: GetContentLength(s),
	}
}
