package main

import "net"

func main() {
	conn, err := net.Dial("tcp", "176.99.160.235:6668")
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("test"))
	conn.Close()
}
