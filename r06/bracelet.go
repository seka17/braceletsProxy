package bracelet

import (
	"bytes"
	"io"
	"net"
	"regexp"

	"github.com/seka17/all/server"
)

var reHeader = regexp.MustCompile(`^#@H[0-9]{2}@#$`)

type Bracelet struct {
	imei    string
	prefix  string
	version string
	conn    net.Conn

	server *server.Server
}

func (this *Bracelet) Init(conn net.Conn, server *server.Server) {
	this.conn = conn
	this.server = server
	this.version = "r06"
	return
}

func (this Bracelet) Version() string {
	return this.version
}

// Read reads all data comming from tracker
func (this Bracelet) Read() ([]byte, error) {
	// Messages can be big
	tmp := make([]byte, 1100)
	if _, err := this.conn.Read(tmp); err != nil {
		return nil, err
	} else {
		array := bytes.Split(tmp, []byte(";"))
		// without verification code and ending bytes
		return bytes.Join(array[:len(array)-2], []byte(";")), nil
	}
}

// Read reads all data comming from tracker
func (this Bracelet) ReadFromSrc(source io.Reader) ([]byte, error) {
	// Messages can be big
	tmp := make([]byte, 1100)
	if _, err := source.Read(tmp); err != nil {
		return nil, err
	} else {
		array := bytes.Split(tmp, []byte(";"))
		// without verification code and ending bytes
		return bytes.Join(array[:len(array)-2], []byte(";")), nil
	}
}

// Read reads all data comming from tracker
func (this Bracelet) Write(msg []byte) (n int, err error) {
	return this.conn.Write(msg)
}

func (this Bracelet) Id() string {
	return this.imei
}

func (this Bracelet) IsMessageFromMe(command []byte) bool {
	header := bytes.SplitN(command, []byte(";"), 2)[0]
	return reHeader.Match(header)
}
