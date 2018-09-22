package ftpparser

import "github.com/manuporto/distributedHTTPServer/pkg/httpparser"

const (
	methodGET    = iota
	methodPOST   = iota
	methodPUT    = iota
	methodDELETE = iota
	methodRESP   = iota

	statusOK          = 200
	statusNotFound    = 404
	statusFileExists  = 410 //Review
	statusServerError = 500
)

var methods = map[string]int{
	"GET":    methodGET,
	"POST":   methodPOST,
	"PUT":    methodPUT,
	"DELETE": methodDELETE,
}

type FTPPacket struct {
	method int
	path   string
	body   []byte
}

type FTPResponse struct {
	status int
	svMsg  string
	path   string
	body   []byte
}

func HTTPToFTP(hf *httpparser.HttpFrame) *FTPPacket {
	return &FTPPacket{methods[hf.Header.Method], hf.Header.URI, hf.Body}
}
