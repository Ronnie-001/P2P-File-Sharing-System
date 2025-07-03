package discovery

import (
	"p2p-file-share/internals/ui"

	"encoding/base64"
	"encoding/pem"

	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	"context"
	"fmt"
	"time"

	"github.com/grandcat/zeroconf"
)

func StartServer(identity string) (*zeroconf.Server, error) {

	fmt.Println("-> Starting the mDNS server.")
	
	privateKeyBlock, publicKeyBlock := GetEncodedRsaKeys()
	
	privateKeyData, _ := pem.Decode(privateKeyBlock)

	if privateKeyData == nil {
		return nil, fmt.Errorf("no PEM data found for privateKeyBlock")
	} else if privateKeyData.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid headers for private key (PEM formatted block)")	
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyData.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key from PEM format: %v", err)
	}

	hash := sha256.Sum256([]byte(identity))
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash[:]) 
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt signature: %v", err)
	}
	
	publicKeyBlockStr := base64.StdEncoding.EncodeToString(publicKeyBlock)

	server, err := zeroconf.Register(
		"P2P fileshare",
		"_fileshare._tcp",
		".local",
		8000,
		// TODO: Inlucde a signed payload within the TXT field, should contain
		// the users actual name, the pulbicKeyBlock & a timestamp. Also, encode this to 
		// base64

		[]string{
			"A simple file sharing service.",  
			string(signature), 
			publicKeyBlockStr,
			},
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
	go func (results <-chan *zeroconf.ServiceEntry) () {
		for entry := range entries {
			// TODO: Decode and parse payload for the publicKeyBytes

			// grab the public key & signature from the txt field
			signature := []byte(entry.Text[1])
			publicKeyBlockStr := entry.Text[2]
			
			publicKeyBlock, err := base64.StdEncoding.DecodeString(publicKeyBlockStr)
			if err != nil {
				fmt.Println("error when decoding PEM block string")
			}
			
			fmt.Println(string(publicKeyBlock))
			
			publicKeyData, _ := pem.Decode(publicKeyBlock)
			if publicKeyData == nil {
				fmt.Printf("unable to decode the public key block from mDNS broadcast. ")
			} else if publicKeyData.Type != "RSA PUBLIC KEY"{
				fmt.Printf("invalid headers for public key")
			}

			publicKey, err := x509.ParsePKCS1PublicKey(publicKeyData.Bytes)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	
	err = resolver.Browse(ctx, "_fileshare._tcp", "local", entries)
	if err != nil {
		return fmt.Errorf("error when browsing for mDNS servers: %v", err)
	}
	
	<-ctx.Done()
	return nil
}
