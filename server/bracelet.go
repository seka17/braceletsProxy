package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"reflect"

	"github.com/seka17/all/structures"
)

type Bracelet interface {
	Init(net.Conn, *Server)
	Read() ([]byte, error)
	ReadFromSrc(io.Reader) ([]byte, error)
	Write([]byte) (int, error)
	Id() string
	Parse([]byte) ([]byte, error)
	Version() string

	IsMessageFromMe([]byte) bool

	// Configure bracelet
	Configure(map[string]interface{}) error

	// Commands
	Call(string) error
	PrintMessage(string) error
	Reset() error
	PowerOff() error
	SetContacts([]structures.Contact) error
}

func (this *Server) AddSupportedBracelets(bracelets ...interface{}) {
	for _, v := range bracelets {
		if _, ok := v.(Bracelet); !ok {
			continue
		}
		// Get type of bracelet
		t := reflect.TypeOf(v).Elem()
		this.supportedBracelets = append(this.supportedBracelets, t)
	}
}

// func (this *Server) initBracelet() Bracelet {
// 	return reflect.New(this.supportedBracelets[0]).Interface().(Bracelet)
// }

func (this Server) chooseBracelet(conn net.Conn) (*bytes.Buffer, Bracelet, error) {
	// Large buffer
	msg := make([]byte, 2000)
	_, err := conn.Read(msg)
	if err != nil {
		return nil, nil, err
	}

	for _, t := range this.supportedBracelets {
		br := reflect.New(t).Interface().(Bracelet)
		if br.IsMessageFromMe(msg) {
			buf := &bytes.Buffer{}
			if _, err := buf.Write(msg); err != nil {
				return nil, nil, err
			}
			return buf, br, nil
		}
	}

	return nil, nil, fmt.Errorf("Message not understood by bracelets -> %s", string(msg))

}
