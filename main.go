package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"p2p-file-share/internals/discovery"
	"p2p-file-share/internals/transfer"
	"p2p-file-share/internals/ui"
)

func main() {
	// Grab the identity of the user from command line arguments.
	name :=	ui.GetIdentity()
	
	server, err := discovery.RegisterMDNS(name)
	if err != nil {
		log.Fatalf("Error when starting mDNS server: %v", err)
	}

	fmt.Printf("-> mDNS server started!\n")
	
	users, err := discovery.DiscoverMDNS()
	if err != nil {
		log.Fatalf("error when finding client using service: %v", err)
	}
	
	defer discovery.StopMDNS(server)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("User: %v, connected \n", name)
	

	go transfer.StartTCPServer()	

	var wg sync.WaitGroup 
	wg.Add(1)
	
	/*
		Users should be able to provide user input whilst TCP server listens for connections,
		Run TCP server concurrently with main for loop
	*/
	go func() {
		for {
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			
			input := strings.TrimSpace(text) 
			
			if strings.Compare(input, "exit") == 0 {
				wg.Done()
			}
			
			if strings.Compare(input, "show user list") == 0 {
				ui.DisplayUsers(users)
			}
		}
	}()
	
	wg.Wait()
}
