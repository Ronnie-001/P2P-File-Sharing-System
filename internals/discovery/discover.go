package discovery

import (
	"fmt"
	"github.com/hashicorp/mdns"
)

func StartServer() {
	instanceName := "File sharing service"	
	info := []string{"A peer-to-peer file sharing service."}

	service, _ := mdns.NewMDNSService(instanceName, "_fileshare._tcp.", "", "", 8000, nil, info)
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()
	
	entriesChan := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesChan {
			fmt.Printf("entry: %v\n", entry)
		}
	}()

	mdns.Lookup("_fileshare._tcp.", entriesChan)
	close(entriesChan)
}
