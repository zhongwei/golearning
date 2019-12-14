package main

import (
	"io/ioutil"
	"fmt"
	"crypto/x509"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
)

const privateKeyFile = "./privateKey.pem"
const publicKeyFile = "./publicKey.pem"

func rsaPubEncrypt(fileName string, plainText []byte) ([]byte, error) {
	info, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(info)

	derText := block.Bytes
	publicKey , err := x509.ParsePKCS1PublicKey(derText)

	if err != nil {
		return nil, err
	}

	cipherData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)

	if err != nil {
		return nil, err
	}

	return cipherData, nil
}

func rsaPriDecrypt(fileName string, cipherData []byte) ([]byte, error) {
	info, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(info)

	derText := block.Bytes
	privateKey , err := x509.ParsePKCS1PrivateKey(derText)

	if err != nil {
		return nil, err
	}

	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherData)

	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func main()  {
	src := []byte("Go away!")	
	fmt.Printf("src = %s\n", src)

	cipherData, err := rsaPubEncrypt(publicKeyFile, src)
	if err != nil {
		fmt.Println("Public key encryption failed!")
	}

	fmt.Printf("Cipher data : %x\n", cipherData)

	plainText, err := rsaPriDecrypt(privateKeyFile, cipherData)
	if err != nil {
		fmt.Println("Private key decryption failed!")
	}

	fmt.Printf("plainText data : %s\n", plainText)
}