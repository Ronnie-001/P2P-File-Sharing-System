package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"p2p-file-share/internals/discovery"
	"p2p-file-share/internals/ui"
)

func main() {
	// Grab the identity of the user from command line arguments.
	name :=	ui.SetIdentity()
	server, err := discovery.StartServer()
	if err != nil {
		log.Fatalf("Error when starting mDNS server: %v", err)
	}
	
	defer discovery.StopServer(server)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("User: %v, connected \n", name)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		
		input := strings.TrimSpace(text) 
		
		if strings.Compare(input, "exit") == 0 {
			break
		}
	}
}
