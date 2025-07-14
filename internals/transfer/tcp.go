// Package transfer: used for sending and reciving files between clients.
package transfer

import (
	"fmt"
	"io"
	"log"
	"net"

	"os"
)

var (
	network = "tcp"
	port = ":4500"
	byteLimit = 15000000 

	// map for users and their local IP's
	m = make(map[string]string)
)

func AddIP(name, ip string) {
	m[name] = ip		
}

func StartTCPServer() (conn net.Conn) {
	ln, err := net.Listen(network, port)
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

func SendFile(name, path string) {
	localIP, ok := m[name]
	if !ok {
		fmt.Println("IP of user " + name + " not found!")
	}
	
	address := localIP + port
	fmt.Println(address)
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path) 
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	
	if _, err := io.Copy(conn, file); err != nil {
		log.Fatal(err)
	}
}

func RecieveFile(conn net.Conn) {
	defer conn.Close()

	// create buffer to store our data.
	buffer := make([]byte, byteLimit)

	file, err := os.Create("test")
		if err != nil {
			log.Fatal(err)
		}

	// Read in the data from the sender.
	for {
		_, err := conn.Read(buffer)
		
		// check if we have reached the end of the file
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
	}
	
	file.Write(buffer)
}
