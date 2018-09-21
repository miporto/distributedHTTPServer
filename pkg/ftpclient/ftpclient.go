package ftpclient

import (
	"encoding/gob"
	"net"

	"github.com/manuporto/distributedHTTPServer/pkg/ftpparser"
)

type FTPClient struct {
	conn net.Conn
}

func Connect(address string) (*FTPClient, error) {
	conn, err := net.Dial("tcp4", address)
	if err != nil {
		return nil, err
	}
	return &FTPClient{conn}, nil
}

func (fc FTPClient) Send(p ftpparser.FTPPacket) {
	encoder := gob.NewEncoder(fc.conn)
	encoder.Encode(p)
}

func (fc FTPClient) Receive() ftpparser.FTPResponse {
	decoder := gob.NewDecoder(fc.conn)
	res := &ftpparser.FTPResponse{}
	decoder.Decode(res)
	return *res
}

func (fc FTPClient) Close() {
	fc.conn.Close()
}
