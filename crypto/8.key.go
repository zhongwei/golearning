package main

import (
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"crypto/rand"
	"fmt"
	"os"
)

const privateKeyFile = "./privateKey.pem"
const publicKeyFile = "./publicKey.pem"

func generateKeyPair(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {
		return err
	}

	privateDerText := x509.MarshalPKCS1PrivateKey(privateKey)

	block := pem.Block {
		Type: "SH RSA PRIVATE KEY",
		Headers: nil,
		Bytes: privateDerText,
	}

	fileHandler, err := os.Create(privateKeyFile)
	if err != nil {
		return err
	}

	defer fileHandler.Close()

	err = pem.Encode(fileHandler, &block)

	if err != nil {
		return err
	}

	publicKey := privateKey.PublicKey
	publicDerText := x509.MarshalPKCS1PublicKey(&publicKey)

	block = pem.Block {
		Type: "SH RSA PUBLIC KEY",
		Headers: nil,
		Bytes: publicDerText,
	}

	fileHandler2, err := os.Create(publicKeyFile)
	if err != nil {
		return err
	}

	defer fileHandler2.Close()

	err = pem.Encode(fileHandler2, &block)  

	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Printf("Generate rsa key pair ... \n")
	err := generateKeyPair(2048)
	if err != nil {
		fmt.Printf("Generate rsa key pair failed, err: v%", err)
	}
}