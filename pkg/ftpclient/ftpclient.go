package ftpclient

import (
	"encoding/gob"
	"net"

	"github.com/manuporto/distributedHTTPServer/pkg/ftpparser"
)

type FTPClient struct {
	conn net.Conn
}

func (fc FTPClient) send(p ftpparser.FTPPacket) {
	encoder := gob.NewEncoder(fc.conn)
	encoder.Encode(p)
}

func (fc FTPClient) receive() ftpparser.FTPResponse {
	decoder := gob.NewDecoder(fc.conn)
	res := &ftpparser.FTPResponse{}
	decoder.Decode(res)
	return *res
}
