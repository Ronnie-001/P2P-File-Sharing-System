package discovery

import (
	"github.com/hashicorp/mdns"
)

func StartServer() mdns.Server {
	instanceName := "File sharing service"	
	info := []string{"A peer-to-peer file sharing service."}

	service, _ := mdns.NewMDNSService(instanceName, "_fileshare._tcp.", "", "", 8000, nil, info)
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})

	return *server
}

func StopServer(server *mdns.Server) {
	server.Shutdown()
} 
