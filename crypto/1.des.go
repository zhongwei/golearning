package main

import (
	"bytes"
	"crypto/des"
	"crypto/cipher"
	"fmt"
)

func desCBCEncrypt(src, key []byte) []byte {
	block, err := des.NewCipher(key)
	fmt.Printf("Cipher Block Chaining block size : %d\n", block.BlockSize())

	if err != nil {
		panic(err)
	}

	src = paddingInfo(src, block.BlockSize())
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)

	blockMode.CryptBlocks(src, src)

	return src
}

func desCBCDecrypt(cipherData, key []byte) []byte {
	block, err := des.NewCipher(key)
	fmt.Printf("Cipher Block Chaining block size : %d\n", block.BlockSize())

	if err != nil {
		panic(err)
	}

	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	blockMode := cipher.NewCBCDecrypter(block, iv)

	blockMode.CryptBlocks(cipherData, cipherData)

	plainText := unpaddingInfo(cipherData)
	return plainText

}

func paddingInfo(src []byte, blockSize int) []byte {
	length := len(src)

	remain := length % blockSize
	paddingNumber := blockSize - remain

	s1 := byte(paddingNumber)

	s2 := bytes.Repeat([]byte{s1}, paddingNumber)

	srcNew := append(src, s2...)

	return srcNew
}

func unpaddingInfo(plainText []byte) []byte {
	length := len(plainText)

	if length == 0 {
		return []byte{}
	}

	lastByte := plainText[length-1]

	unpaddingNumber := int(lastByte)

	return plainText[:length - unpaddingNumber]
}

func main() {
	src := []byte("123456789")
	key := []byte("12345678")

	fmt.Printf("src : %s\n", src)

	cipherData := desCBCEncrypt(src, key)
	fmt.Printf("cipherData : %x\n", cipherData)

	plainText := desCBCDecrypt(cipherData, key)
	fmt.Printf("plainText : %s\n", plainText)
}