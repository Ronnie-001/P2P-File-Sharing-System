package discovery

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"encoding/pem"

	"fmt"
)

func generatePrivateKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Errorf("error when generating private key")  
	}

	publicKey := &privateKey.PublicKey
	
	return privateKey, publicKey
}

func encodeRsaKeys() (privateKeyPEM, publicKeyPEM *pem.Block) {
	
	private, public := generatePrivateKey()	

	privateKeyPEM = &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	}

	publicKeyPEM = &pem.Block{
		Type: "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(public),
	}
	
	return
}
