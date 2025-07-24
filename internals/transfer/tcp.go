// Package transfer: used for sending and reciving files between clients.
package transfer

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
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

func StartTCPServer() (error) {
	ln, err := net.Listen(network, port)
	if err != nil {
		log.Fatal(err)
	}
	
	data := make(chan string)
	var wg sync.WaitGroup
	
	// only accept 2 connections 1=[name & filetype] & 2=[actual data from the file]  
	for i := 0; i < 2; i++ {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("unable to accept incoming connections: %v", err)		
		}
		wg.Add(1)
		go handleConnection(i, conn, &wg, data)
	}

	wg.Wait()
	close(data)
	RecieveFile(data)	

	return nil
}

func handleConnection(id int, conn net.Conn, wg *sync.WaitGroup, data chan string) (error) {
	defer conn.Close()
	defer wg.Done()
	
	bytes := make([]byte, byteLimit)	
	for {
		_, err := conn.Read(bytes)	
		if err != nil {
			if id == 0 {
				return fmt.Errorf("error when trying to read in the name and filetype: %v", err)
			}
			if id == 1 {
				return fmt.Errorf("error when trying to read in raw data from file")	
			}
		}

		data <- string(bytes)
		fmt.Printf("Data read from %v successfully!", id)
	}
}

func SendFile(name, path string) {
	localIP, ok := m[name]
	if !ok {
		fmt.Println("IP of user " + name + " not found!")
	}
	
	// form the address from the local ip and port number
	address := localIP + port
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

func RecieveFile(data chan string) {
	//TODO: iterate over the go channel to recive file info and data.
}
