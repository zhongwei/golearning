package main

import (
	"errors"
	"math/big"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

const ECCPrivateKeyFile = "./ECCPrivateKey.pem"
const ECCPublicKeyFile = "./ECCPublicKey.pem"

func getPublicKey(fileName string) (*ecdsa.PublicKey, error) {
	info, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(info)
	derText := block.Bytes

	keyInterface, err := x509.ParsePKIXPublicKey(derText)

	if err!= nil {
		fmt.Println("Parse public key failed!")
		return nil, err
	}

	publicKey, OK := keyInterface.(*ecdsa.PublicKey)
	if !OK {
		return nil, errors.New("Reflect public key failed!")
	}

	return publicKey, nil
}

func getPrivateKey(fileName string) (*ecdsa.PrivateKey, error) {
	info, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(info)
	derText := block.Bytes

	return x509.ParseECPrivateKey(derText)

}

type Signature struct {
	r, s *big.Int
}

func sign(src []byte) (Signature, error) {
	privateKey, err := getPrivateKey(ECCPrivateKeyFile)
	if err != nil {
		fmt.Println("Get private key failed!")
	}

	hash := sha256.Sum256(src)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		fmt.Println("Sign failed!")
	}

	return Signature{r, s}, nil
}

func verify(src []byte, sig Signature) bool {
	publicKey, err := getPublicKey(ECCPublicKeyFile)
	if err != nil {
		fmt.Println("Get public key failed!")
	}

	hash := sha256.Sum256(src)

	return ecdsa.Verify(publicKey, hash[:], sig.r, sig.s)
}

func main() {
	src := []byte("Go away!")
	fmt.Printf("src = %s\n", src)

	signature, err := sign(src)
	if err != nil {
		fmt.Println("Sign failed！")
	}
	
	fmt.Printf("Signed : %x\n", signature)

	isValid := verify(src, signature)

	if isValid {
		fmt.Println("Verify success！")
		return 
	}

	fmt.Println("Verify failed！")
}
