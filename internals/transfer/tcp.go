// Package transfer: used for sending and reciving files between clients.
package transfer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	network = "tcp"
	port = ":4500"

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
			return fmt.Errorf("Unable to accept incoming connection: %v", err)		
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
		log.Printf("Unable to grab metadata from incoming connection")
		return 
	}

	metadataStr = strings.TrimSpace(metadataStr)
	parts := strings.Split(metadataStr, "|")
	if len(parts) != 2 {
		log.Printf("Metadata does not have all avaliable parts")
		return
	}

	filename := parts[0]
	fileSizeStr := parts[1]
	fileSize, err :=  strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		log.Printf("Unable to parse metadata for file size")
		return
	}

	// create the file
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Unable to create file")
		return
	}

	defer file.Close()
	
	n, err := io.Copy(file, io.LimitReader(reader, fileSize))
	if err != nil {
		log.Printf("Unable to copy file contents to TCP server")
		return 
	}

	if n != fileSize {
		log.Printf("Recieved more bytes than usual")
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
		log.Printf("Failed to establish connection")
	}
	
	defer conn.Close()
	
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal("Unable to get file info")
	}

	filename := fileInfo.Name()
	fileSize := fileInfo.Size()

	metadata := fmt.Sprintf("%s|%d\n", filename, fileSize)
	_, err = conn.Write([]byte(metadata))
	if err != nil {
		log.Printf("Unable to send over metadata")
	}

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to grab file from path:" + path)
	}

	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal("Unable to send file over TCP connection")
	}
}
