package discovery

import (
	"fmt"

	"github.com/grandcat/zeroconf"
)

func StartServer() (*zeroconf.Server, error) {
	fmt.Println("-> Starting the mDNS server.")
	
	server, err := zeroconf.Register(
		"P2P fileshare",
		"_fileshare._tcp",
		".local",
		8000,
		[]string{"desc=A simple file sharing service."},
		nil,
	) 
	if err != nil {
		return nil, fmt.Errorf("couldn't register the mDNS server: %v", err) 
	}

	return server, nil
}

func DiscoverServers() {

}

func StopServer(server *zeroconf.Server) {
	server.Shutdown()	
}
