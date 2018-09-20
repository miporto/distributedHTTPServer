package ftpparser

const (
	methodGET = iota
	methodPOST = iota
	methodPUT = iota
	methodDELETE = iota
	methodRESP = iota

	statusOK = 200
	statusNotFound = 404
	statusFileExists = 410 //Review
	statusServerError = 500

)

type FTPPacket struct {
	method int
	pathLen int
	path string
	bodyLen int
	body []byte
}

type FTPResponse struct {
	status int
	svMsgLen int
	svMsg string
	pathLen int
	path string
	bodyLen int
	body []byte
}
