package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

type Server struct {
	conn      *net.TCPListener
	Bracelets map[string]Bracelet
	logger    *logrus.Logger
	sync.RWMutex
}

type Bracelet interface {
	Read() ([]byte, error)
	Id() string
}

type R06 struct {
	imei   string
	conn   net.Conn
	server *Server
}

func initLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.StandardLogger()
	logger.Level = level
	logger.Formatter = &logrus.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
	}
	return logger
}

func main() {

	logger := initLogger(logrus.DebugLevel)

	ad := flag.String("gate", ":6669", "Address of gate server")
	service := flag.String("service", ":6670", "Address of service server")
	flag.Parse()

	parsedService := strings.Split(*service, ":")

	addr, err := net.ResolveTCPAddr("tcp", *ad)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	logger.Info("Listening for connections...")

	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		logger.Debug("Got connection")
		if err != nil {
			logger.WithError(err).Error("Accept connection")
			continue
		}
		tmp := make([]byte, 1100)
		if _, err := conn.Read(tmp); err != nil {
			logger.Error(err)
			continue
		}
		logger.WithField("msg", string(tmp)).Debug("Request")

		if _, err := conn.Write(append([]byte(fmt.Sprintf("ServiceIP:%s,%s;", parsedService[0], parsedService[1])), []byte{0x010, 0x010, 0x01}...)); err != nil {
			logger.WithError(err).Error("Write")
			continue
		}
		logger.WithField("msg", string(append([]byte(fmt.Sprintf("ServiceIP:%s,%s;", parsedService[0], parsedService[1])), []byte{0x010, 0x010, 0x01}...))).Debug("Response")

		time.Sleep(time.Second)
		conn.Close()
	}
}
