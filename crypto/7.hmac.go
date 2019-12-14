package main

import (
	"crypto/hmac"
	"fmt"
	"crypto/sha256"
)

func genHMAC(src []byte, key []byte) []byte {
	hasher := hmac.New(sha256.New, key)
	hasher.Write(src)

	return hasher.Sum(nil)
}

func verifyHMAC(src, key, mac []byte) bool {
	return hmac.Equal(genHMAC(src, key), mac)
}

func main() {
	src := []byte("hello world")
	key := []byte("1234567890")

	mac1 := genHMAC(src, key)
	fmt.Printf("mac1 : %x\n", mac1)
	fmt.Printf("Is equal? : %v\n", verifyHMAC(src, key, mac1))
	fmt.Printf("Is equal? : %v\n", verifyHMAC([]byte("hello world!"), key, mac1))
}