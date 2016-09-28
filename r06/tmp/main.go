package main

import (
	"bytes"
	"flag"
	"net"
	"sync"

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

	ad := flag.String("addr", ":6669", "Address")
	flag.Parse()

	server := &Server{Bracelets: make(map[string]Bracelet), logger: initLogger(logrus.DebugLevel)}
	addr, err := net.ResolveTCPAddr("tcp", *ad)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	server.logger.Info("Listening for connections...")

	server.conn = listener
	defer server.conn.Close()

	for {
		conn, err := server.conn.AcceptTCP()
		server.logger.Debug("Got connection")
		if err != nil {
			server.logger.WithError(err).Error("Accept connection")
			continue
		}
		go server.registrateBracelet(conn)
	}
}

func (this *Server) registrateBracelet(conn net.Conn) {
	br := &R06{
		conn: conn,
		// imei: "12345",
	}
	defer conn.Close()

	this.logger.Debug("Bracelet connected!")

	// this.AddBracelet(br)

	for {
		msg, err := br.Read()
		if err != nil {
			this.logger.Debug("Bracelet disconnected!")
			return
		}

		this.logger.WithField("msg", string(msg)).Info("Got message")
		this.logger.Debug(msg)

		// if _, err := br.Write(append([]byte("ServiceIP:176.99.160.235,6670;"), []byte{0x010, 0x010, 0x01}...)); err != nil {
		// 	this.logger.WithError(err).Error("Write")
		// 	return
		// }

		// this.logger.WithField("msg", string(append([]byte("ServiceIP:176.99.163.176,6668;"), []byte{0x010, 0x010, 0x01}...))).Info("Send message")

	}

	// br.BracConn.SetDeadline(time.Now().Add(DeadLineTime))

	// br.ID, br.Version = br.ParseRegistrationRequest(req)
	// if br.ID == "" {
	// 	server.Logger.Debug("Bad registration request", "command", string(req))
	// 	return
	// }
	// fmt.Println(">> Have connection from bracelet", br.ID)
	// server.Logger.Info("Request", "IMEI", br.ID, "command", string(req))
	// br.Logger = server.Logger.New("IMEI", br.ID)

	// fmt.Println(">> Before add bracelet", br.ID)
	// server.AddBracelet(br)
	// fmt.Println(">> After add bracelet", br.ID)

	// answer := br.AnswerRegistrationRequest()
	// br.Logger.Info("Response", "command", answer)

	// if _, err := br.Write(answer); err != nil {
	// 	br.Logger.Warn("Write", "error", err)
	// 	return
	// }
	// fmt.Println(">> Before registrate in fs bracelett", br.ID)
	// br.RegistrateInFs()

	// br.RLock()
	// fmt.Println(">> After registrate in fs bracelett", br.ID, br.Registred)
	// registred := br.Registred
	// br.RUnlock()
	// if registred {
	// 	br.Push()
	// 	go func() { time.Sleep(time.Second * 4); br.SetPhones() }()
	// }

	// fmt.Println(">> Listen bracelett", br.ID)
	// if err := server.ListenBracelet(br.ID); err != nil {
	// 	br.Logger.Warn("Connection closed", "error", err)
	// }
	// return
}

func (this *Server) AddBracelet(br Bracelet) {
	this.Lock()
	defer this.Unlock()
	this.Bracelets[br.Id()] = br
	this.logger.WithField("bracelet", br.Id()).Info("Adding bracelet")
	return
}

func (this Server) GetBracelet(id string) (Bracelet, bool) {
	this.RLock()
	br, ok := this.Bracelets[id]
	this.logger.WithFields(logrus.Fields{
		"bracelet": br.Id(),
		"ok":       ok,
	}).Debug("Getting bracelet from map")
	this.RUnlock()
	return br, ok
}

func (this *Server) RemoveBracelet(id string) {
	this.Lock()
	br, ok := this.Bracelets[id]
	delete(this.Bracelets, id)
	this.Unlock()
	if ok {
		this.logger.WithField("bracelet", br.Id()).Info("Removing bracelet")
	}
	return
}

// Read reads all data comming from tracker
func (this R06) Read() ([]byte, error) {
	// Messages can be big
	tmp := make([]byte, 1100)
	if _, err := this.conn.Read(tmp); err != nil {
		return nil, err
	} else {
		array := bytes.Split(tmp, []byte(";"))
		// without verification code and ending bytes
		return bytes.Join(array[:len(array)-2], ";"), nil
	}
}

// Read reads all data comming from tracker
func (this R06) Write(msg []byte) (n int, err error) {
	return this.conn.Write(msg)
}

func (this R06) Id() string {
	return this.imei
}

var bufferPool = &sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func initBuffers() {
	for i := 0; i < 300; i++ {
		bufferPool.Put(new(bytes.Buffer))
	}
}

// //ListenBracelets listens bracelets and send the further
// func (server *Server) ListenBracelets() {
// 	fmt.Println("In listen Bracelets!")
// 	for {
// 		conn, err := server.BC.AcceptTCP()
// 		fmt.Println("Got connection")
// 		if err != nil {
// 			server.Logger.Error("Accept connection", "error", err)
// 			continue
// 		}
// 		go server.RegistrateBracelet(conn)
// 	}
// }

// // ListenBracelet listens bracelet and sends obtained requests further
// func (server *Server) ListenBracelet(id string) error {
// 	br, ok := server.GetBracelet(id)
// 	if !ok {
// 		server.Logger.Warn("There's no such bracelet", "IMEI", id)
// 		return nil
// 	}

// 	for {
// 		request, err := br.Read()
// 		if err != nil {
// 			return err
// 		}
// 		go func() {
// 			response, err := server.AnswerRequest(id, string(request))
// 			switch err {
// 			case nil:
// 				if _, err = br.Write(response); err != nil {
// 					br.BracConn.Close()
// 				}
// 				return
// 			case NoResponse:
// 				return
// 			default:
// 				br.Logger.Warn("HandleRequest", "error", err)
// 			}
// 		}()
// 	}
// }
