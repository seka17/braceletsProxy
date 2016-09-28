package bracelet

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"

	"github.com/seka17/all/server"
)

var reHeader = regexp.MustCompile(`^\[[A-Z0-9]{2}$`)
var reImei = regexp.MustCompile(`^[0-9]{10}$`)

type Bracelet struct {
	imei    string
	prefix  string
	version string
	conn    net.Conn
	server  *server.Server
}

func (this *Bracelet) Init(conn net.Conn, server *server.Server) {
	this.conn = conn
	this.server = server
	this.version = "q80"
	return
}

func (this Bracelet) Version() string {
	return this.version
}

// Read reads all data comming from tracker
// Example: [3G*6005060193*0009*LK,0,0,57]
func (this *Bracelet) Read() ([]byte, error) {
	tmp := make([]byte, 20)
	// Read header
	if _, err := this.conn.Read(tmp); err != nil {
		return nil, err
	}
	fmt.Println("Header of read ->", string(tmp))

	// Set imei from header
	if this.imei == "" {
		this.imei = string(tmp[4:14])
		this.prefix = string(tmp[1:3])
		this.server.AddBracelet(this)
	} // Get length of body

	length, err := strconv.ParseInt(string(tmp[16:19]), 16, 10)
	if err != nil {
		return nil, err
	}

	// +1 for last "]"
	tmp = make([]byte, length+1)
	// Read header
	if _, err := this.conn.Read(tmp); err != nil {
		return nil, err
	} else {
		return tmp[:len(tmp)-1], nil
	}
}

// Read reads all data comming from source
func (this *Bracelet) ReadFromSrc(source io.Reader) ([]byte, error) {
	tmp := make([]byte, 20)
	// Read header
	if _, err := source.Read(tmp); err != nil {
		return nil, err
	}
	fmt.Println("Header of read ->", string(tmp))

	// Set imei from header
	if this.imei == "" {
		this.imei = string(tmp[4:14])
		this.prefix = string(tmp[1:3])
		this.server.AddBracelet(this)
	}
	// Get length of body

	length, err := strconv.ParseInt(string(tmp[16:19]), 16, 10)
	if err != nil {
		return nil, err
	}

	// +1 for last "]"
	tmp = make([]byte, length+1)
	// Read header
	if _, err := source.Read(tmp); err != nil {
		return nil, err
	} else {
		return tmp[:len(tmp)-1], nil
	}
}

// Read reads all data comming from tracker
func (this Bracelet) Write(msg []byte) (n int, err error) {
	if msg == nil {
		return 0, nil
	}
	return this.conn.Write(msg)
}

func (this Bracelet) Id() string {
	return this.imei
}

func (this Bracelet) AddHeader(msg string) []byte {
	tmp := strconv.FormatInt(int64(len(msg)), 16)
	appendix := "0000"
	tmp = appendix[len(tmp):] + tmp
	return []byte(fmt.Sprintf("[%s*%s*%s*%s]", this.prefix, this.imei, tmp, msg))
}

func (this Bracelet) IsMessageFromMe(command []byte) bool {
	tmp := bytes.SplitN(command, []byte("*"), 3)
	return reHeader.Match(tmp[0]) && reImei.Match(tmp[1])
}
