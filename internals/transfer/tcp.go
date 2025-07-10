package transfer

import (
	"log"
	"net"
)

func ListenAndAccept() {
	ln, err := net.Listen("tcp", ":21") // port 21 for file transfers using TCP
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		conn, err := ln.Accept()	
		if err != nil {
			log.Fatal(err)
		}

	}
}
