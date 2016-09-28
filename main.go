package main

import (
	"log"

	q80 "github.com/seka17/all/q80"
	r06 "github.com/seka17/all/r06"
	"github.com/seka17/all/server"
)

func main() {
	config := server.Config{
		Address:     ":6668",
		AddressName: "176.99.176.255",
	}

	s := server.Init(config)

	s.AddSupportedBracelets(&q80.Bracelet{}, &r06.Bracelet{})

	if err := s.Run(); err != nil {
		log.Println("Error", err)
	}
}
