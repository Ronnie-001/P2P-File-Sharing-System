package discovery

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hashicorp/mdns"
)

func StartServer() (*mdns.Server, error) {
	instanceName := "File sharing service"	
	info := []string{"A peer-to-peer file sharing service."}
	
	var ips []net.IP
	addrs, _ := net.InterfaceAddrs()
	
	// add all NON-LOOPBACK ip's to mDNS config.
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			} 
		}
	}

	service, err := mdns.NewMDNSService(
		instanceName, 
		"_fileshare._tcp",
		"",
		"",
        8000,
        ips,
        info,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new mDNS service: %v", err)	
	}

	// start the mDNS server.
	server, err := mdns.NewServer(&mdns.Config{
		Zone: service,
		Logger: log.New(os.Stdout, "[mDNS DEBUG]", log.LstdFlags),	
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start up the mDNS server: %v", err)
	}

	return server, nil
}

func StopServer(server *mdns.Server) {
	server.Shutdown()
} 
