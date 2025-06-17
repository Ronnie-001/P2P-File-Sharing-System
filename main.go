package main

import (
	"fmt"
	"p2p-file-share/internals/discovery"
)

func main() {
	fmt.Println("Peer-to-Peer file sharing system!")
	discovery.StartServer()	
}
