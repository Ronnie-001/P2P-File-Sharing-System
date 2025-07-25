// Package transfer: used for sending and reciving files between clients.
package transfer

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("unable to accept incoming connections: %v", err)		
		}
		go handleConnection(conn)		
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// grab the files metadata
	reader := bufio.NewReader(conn)
	metadataStr, err := reader.ReadString('\n')
	if err != nil {
		return 
	}

	metadataStr = strings.TrimSpace(metadataStr)
	parts := strings.Split(metadataStr, "|")
	if len(parts) != 2 {
		return
	}

	filename := parts[0]
	fileSizeStr := parts[1]
	fileSize, err :=  strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		return
	}

	// create the file
	file, err := os.Create(filename)
	if err != nil {
		return
	}

	defer file.Close()
	
	n, err := io.Copy(file, io.LimitReader(reader, fileSize))
	if err != nil {
		return 
	}

	if n != fileSize {
		log.Printf("Recieved more bytes than usual")
	}

}

func SendFile(name, path string) {
	var wg sync.WaitGroup

	localIP, ok := m[name]
	if !ok {
		fmt.Println("IP of user " + name + " not found!")
	}
	
	// form the address from the local ip and port number
	address := localIP + port

	// go routine to send the name of the file and the file type.
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := net.Dial(network, address)
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()	

		splitPath := strings.Split(path, "/")
		nameOfFile := splitPath[len(splitPath) - 1]
		
		conn.Write([]byte(nameOfFile))
	}()	
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := net.Dial(network, address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		file, err := os.Open(path) 
		if err != nil {
			log.Fatal(err)
		}	

		defer file.Close()
		if _, err := io.Copy(conn, file); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
