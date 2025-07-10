package discovery

import (
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
	
	privatePEMBytes, publicPEMBytes := GetEncodedRsaKeys()
	
	privateBlock, _ := pem.Decode(privatePEMBytes)
	if privateBlock == nil {
		return nil, fmt.Errorf("no PEM data found for privatePEMBytes")
	} else if privateBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid headers for private key (PEM formatted block)")	
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key from PEM format: %v", err)
	}

	hash := sha256.Sum256([]byte(identity))
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash[:]) 
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt signature: %v", err)
	}

	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	
	publicKeyBlockStr := base64.StdEncoding.EncodeToString(publicPEMBytes)
	
	// split the publicKeyBlock (PEM block) into chunks due to TXT field byte limit. 
	chunk1 := publicKeyBlockStr[:200]
	chunk2 := publicKeyBlockStr[200:]

	server, err := zeroconf.Register(
		"P2P fileshare",
		"_fileshare._tcp",
		".local",
		8000,
		[]string{
			"A simple file sharing service.",  
			encodedSignature,
			chunk1,		// publicKeyBlockStr broken down here due to 255 byte limit of txt fields 
			chunk2,
			identity,
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

func DiscoverServers() ([]string, error) {
	
	var userList []string

	resolver, err := zeroconf.NewResolver(nil)	
	if err != nil {
		return userList, fmt.Errorf("error when starting up discovery: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func (results <-chan *zeroconf.ServiceEntry) () {
		for entry := range entries {

			// grab the public key & signature from the txt field
			signature := entry.Text[1]

			chunk1 := entry.Text[2]
			chunk2 := entry.Text[3]		
			publicKeyBlockStr := chunk1 + chunk2 

			publicPEMBlock, err := base64.StdEncoding.DecodeString(publicKeyBlockStr)
			if err != nil {
				fmt.Println("error when decoding PEM block string")
			}
			
			publicBlock, _ := pem.Decode(publicPEMBlock)
			if publicBlock == nil || publicBlock.Type != "RSA PUBLIC KEY" {
				fmt.Printf("unable to decode the public key block from mDNS broadcast. ")
			}

			publicKey, err := x509.ParsePKCS1PublicKey(publicBlock.Bytes)
			if err != nil {
				fmt.Printf("error parsing public block: %v", err)
			}
			
			identity := entry.Text[4]	
			hash := sha256.Sum256([]byte(identity))

			decodedSignature, err := base64.StdEncoding.DecodeString(signature)
			if err != nil {
				fmt.Println("error decoding signature")
			}
			
			err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], decodedSignature) 
			if err != nil {
				fmt.Printf("Unable to verify signature from mDNS service with public key." + entry.ServiceName() + ": %v ", err)
			}
			
			userList = append(userList, identity)
			
		}

	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	
	err = resolver.Browse(ctx, "_fileshare._tcp", "local", entries)
	if err != nil {
		return userList, fmt.Errorf("error when browsing for mDNS servers: %v", err)
	}
	
	<-ctx.Done()
	return userList, nil
}
