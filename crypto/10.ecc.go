package main

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"encoding/pem"
	"crypto/x509"
	"crypto/rand"
	"fmt"
	"os"
)

const ECCPrivateKeyFile = "./ECCPrivateKey.pem"
const ECCPublicKeyFile = "./ECCPublicKey.pem"

func generateKeyPair() error {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}

	privateDerText, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	block := pem.Block {
		Type: "SH ECC PRIVATE KEY",
		Headers: nil,
		Bytes: privateDerText,
	}

	fileHandler, err := os.Create(ECCPrivateKeyFile)
	if err != nil {
		return err
	}

	defer fileHandler.Close()

	err = pem.Encode(fileHandler, &block)

	if err != nil {
		return err
	}

	publicKey := privateKey.PublicKey
	publicDerText, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}

	block = pem.Block {
		Type: "SH ECC PUBLIC KEY",
		Headers: nil,
		Bytes: publicDerText,
	}

	fileHandler2, err := os.Create(ECCPublicKeyFile)
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
	fmt.Printf("Generate ECC key pair ... \n")
	err := generateKeyPair()
	if err != nil {
		fmt.Printf("Generate ECC key pair failed, err: v%", err)
	}
}