package transfer

import (
	"log"
	"net"
)

func StartTCPServer() (conn net.Conn) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		
		go RecieveFile(conn)
	}
}

func SendFile(path string) {}

func RecieveFile(conn net.Conn) {
	defer conn.Close()
}
