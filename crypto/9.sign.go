package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

const privateKeyFile = "./privateKey.pem"
const publicKeyFile = "./publicKey.pem"

func getPublicKey(fileName string) (*rsa.PublicKey, error) {
	info, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(info)
	derText := block.Bytes

	return x509.ParsePKCS1PublicKey(derText)
}

func getPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	info, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(info)
	derText := block.Bytes

	return x509.ParsePKCS1PrivateKey(derText)

}

func sign(src []byte) ([]byte, error) {
	privateKey, err := getPrivateKey(privateKeyFile)
	if err != nil {
		fmt.Println("Get private key failed!")
	}

	hash := sha256.Sum256(src)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}

func verify(src []byte, sig []byte) error {
	publicKey, err := getPublicKey(publicKeyFile)
	if err != nil {
		fmt.Println("Get public key failed!")
	}

	hash := sha256.Sum256(src)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], sig)
}

func main() {
	src := []byte("Go away!")
	fmt.Printf("src = %s\n", src)

	signature, err := sign(src)
	if err != nil {
		fmt.Println("Sign failed！")
	}
	
	fmt.Printf("Signed : %x\n", signature)

	err = verify(src, signature)
	if err != nil {
		fmt.Println("Verify failed！")
		return 
	}

	fmt.Println("Verify success！")
}
