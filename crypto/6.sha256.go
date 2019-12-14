package main

import (
	"os"
 	"fmt"
 	"crypto/sha256"
 	"io"
)

func main() {
	file, err := os.Open("6.sha256.go")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()

	length, err := io.Copy(hasher, file)

	fmt.Printf("length : %d\n", length)

	hash := hasher.Sum(nil)

	fmt.Printf("hash : %x\n", hash)
}