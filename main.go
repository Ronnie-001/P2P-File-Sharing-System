package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"p2p-file-share/internals/discovery"
	"p2p-file-share/internals/ui"
)

func main() {
	// Grab the identity of the user from command line arguments.
	name :=	ui.SetIdentity()
	server := discovery.StartServer()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("User: %v, connected \n", name)
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		
		input := strings.TrimSpace(text) 
		
		if strings.Compare(input, "exit") == 0 {
			discovery.StopServer(&server)
			break
		}
	}
}
