package main

import (
	"log"

	q80 "github.com/seka17/all/q80"
	"github.com/seka17/all/server"
)

func main() {
	config := server.Config{
		Address: ":6668",
	}

	s := server.Init(config)

	s.AddSupportedBracelets(&q80.Bracelet{})

	if err := s.Run(); err != nil {
		log.Println("Error", err)
	}
}
