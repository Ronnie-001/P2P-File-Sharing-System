package discovery

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"

	"github.com/grandcat/zeroconf"
)


func StartServer(identity string) (*zeroconf.Server, error) {
	fmt.Println("-> Starting the mDNS server.")
	
	privateKey, _ := generateKeyPairs() 

	hash := sha256.Sum256([]byte(identity))
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash[:]) 

	server, err := zeroconf.Register(
		"P2P fileshare",
		"_fileshare._tcp",
		".local",
		8000,
		[]string{"desc=A simple file sharing service. signature= " + string(signature)},
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
