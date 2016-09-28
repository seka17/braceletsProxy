package server

import (
	"net"
	"reflect"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

type Server struct {
	Bracelets          map[string]Bracelet
	supportedBracelets []reflect.Type
	address            string
	addressName        string
	conn               *net.TCPListener
	logger             *logrus.Logger
	sync.RWMutex
}

type Config struct {
	LogLevel    *logrus.Level `json:"log_level"`
	Address     string        `json:"address"`
	AddressName string        `json:"addressName"`
}

func Init(config Config) *Server {
	var logLevel logrus.Level
	if config.LogLevel != nil {
		logLevel = *config.LogLevel
	} else {
		logLevel = logrus.DebugLevel
	}
	return &Server{
		Bracelets:          make(map[string]Bracelet),
		address:            config.Address,
		addressName:        config.AddressName,
		logger:             initLogger(logLevel),
		supportedBracelets: make([]reflect.Type, 0),
	}
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

func (this *Server) Run() error {
	addr, err := net.ResolveTCPAddr("tcp", this.address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	this.logger.Info("Listening for connections on " + this.address)

	this.conn = listener
	defer this.conn.Close()

	for {
		conn, err := this.conn.AcceptTCP()
		this.logger.Debug("Got connection")
		if err != nil {
			this.logger.WithError(err).Error("Accept connection")
			continue
		}
		go this.registrateBracelet(conn)
	}
}

func (this *Server) registrateBracelet(conn net.Conn) {
	defer conn.Close()
	buf, br, err := this.chooseBracelet(conn)
	if err != nil {
		this.logger.WithError(err).Error("Can't init bracelet")
		return
	}
	br.Init(conn, this)

	first := true

	defer this.RemoveBracelet(br.Id())

	this.logger.Debug("Bracelet connected!")

	// k := 1

	for {
		msg := make([]byte, 0)
		if first {
			msg, err = br.ReadFromSrc(buf)
			if err != nil {
				this.logger.WithError(err).Error("Can't read from initial buffer")
				return
			}
			buf.Reset()
			first = false
		} else {
			msg, err = br.Read()
			if err != nil {
				this.logger.Debug("Bracelet disconnected!", err)
				return
			}
		}

		this.logger.WithField("msg", string(msg)).Info("Got message")

		if response, err := br.Parse(msg); err != nil {
			this.logger.WithError(err).Error("Can't parse response")
		} else {
			if _, err := br.Write(response); err != nil {
				this.logger.WithError(err).Error("Can't send response")
			} else {
				this.logger.WithField("msg", string(response)).Info("Send message")
			}
		}

		// if k == 1 {
		// 	k = 2
		// 	if err := br.SetContacts([]structures.Contact{structures.Contact{
		// 		Number: "8 (916) 025-26-46",
		// 		Name:   "Sergey",
		// 	}, structures.Contact{
		// 		Number: "8 (915) 015-46-94",
		// 		Name:   "Amina",
		// 	}}); err != nil {
		// 		this.logger.Error(err)
		// 	}
		// }

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
	_, ok := this.Bracelets[id]
	delete(this.Bracelets, id)
	this.Unlock()
	if ok {
		this.logger.WithField("bracelet", id).Info("Removing bracelet")
	}
	return
}

func (this *Server) GetAddress() (string, string) {
	return this.addressName, strings.Split(this.address, ":")[1]
}

// var bufferPool = &sync.Pool{
// 	New: func() interface{} { return new(bytes.Buffer) },
// }

// func initBuffers() {
// 	for i := 0; i < 300; i++ {
// 		bufferPool.Put(new(bytes.Buffer))
// 	}
// }

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
