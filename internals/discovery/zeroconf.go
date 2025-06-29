package discovery

import (
	"context"
	"encoding/pem"
	"p2p-file-share/internals/ui"

	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	"fmt"
	"time"

	"github.com/grandcat/zeroconf"
)

func StartServer(identity string) (*zeroconf.Server, error) {

	fmt.Println("-> Starting the mDNS server.")
	
	privateKeyBlock, publicKeyBlock := GetEncodedRsaKeys()
	
	//TODO: Try passing in all the data from the PEM block, rather than just the bytes. 
	privateKeyData, _ := pem.Decode(privateKeyBlock.Bytes)
	if privateKeyData == nil {
		return nil, fmt.Errorf("no PEM data found for privateKeyBlock")
	} else if privateKeyData.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid headers for private key (PEM formatted block)")	
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key from PEM format")
	}

	hash := sha256.Sum256([]byte(identity))
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash[:]) 
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt signature: %v", err)
	}
	
	server, err := zeroconf.Register(
		"P2P fileshare",
		"_fileshare._tcp",
		".local",
		8000,
		[]string{"A simple file sharing service.", string(signature), string(publicKeyBlock.Bytes)},
		nil,
	) 
	if err != nil {
		return nil, fmt.Errorf("couldn't register the mDNS server: %v", err) 
	}

	return server, nil
}

func StopServer(server *zeroconf.Server) {
	server.Shutdown()	
}

func DiscoverServers() (error) {
	resolver, err := zeroconf.NewResolver(nil)	
	if err != nil {
		return fmt.Errorf("error when starting up discovery: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func (results <-chan *zeroconf.ServiceEntry) {
		for entry := range entries {
			// grab the public key & signature from the txt field
			publicKeyBytes := []byte(entry.Text[2])
			signature := []byte(entry.Text[1])

			
			publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
			if err != nil {
				fmt.Printf("error parsing public key: %v", err)
			}

			// Get the identity of the user
			Identity := ui.GetIdentity()
			hash := sha256.Sum256([]byte(Identity))

			err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature) 
			if err != nil {
				fmt.Printf("Unable to verify public key from mDNS service " + entry.ServiceName() + ": %v", err)
			}
			
			// TODO: Add to some sort of user list that can accessed through a command.
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	err = resolver.Browse(ctx, "_fileshare._tcp", "local", entries)
	if err != nil {
		return fmt.Errorf("error when browsing for mDNS servers: %v", err)
	}
	
	<-ctx.Done()
	return nil
}
