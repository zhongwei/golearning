package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)


func aesCTREncrypt(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	fmt.Printf("Counter block size : %d\n", block.BlockSize())

	if err != nil {
		panic(err)
	}

	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	stream.XORKeyStream(src, src)

	return src
}

func aesCTRDecrypt(cipherData, key []byte) []byte {
 	return aesCTREncrypt(cipherData, key) 
}


func main() {
	src := []byte("123456789")
	key := []byte("1234567887654321")

	fmt.Printf("src : %s\n", src)

	cipherData := aesCTREncrypt(src, key)
	fmt.Printf("cipherData : %x\n", cipherData)

	plainText := aesCTRDecrypt(cipherData, key)
	fmt.Printf("plainText : %s\n", plainText)
}