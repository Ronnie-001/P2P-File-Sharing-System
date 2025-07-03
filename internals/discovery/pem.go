package discovery

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"encoding/pem"
	"log"
)

func GenerateKeyPairs() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey
	
	return privateKey, publicKey
}

func GetEncodedRsaKeys() (privateKeyBlockBytes []byte, publicKeyBlockBytes []byte) {
	
	private, public := GenerateKeyPairs()	
	
	privateKeyBlock := &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	}

	privateKeyBlockBytes = pem.EncodeToMemory(privateKeyBlock)
	
	publicKeyBlock := &pem.Block{
		Type: "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(public),
	}
	
	publicKeyBlockBytes = pem.EncodeToMemory(publicKeyBlock)

	return
}
