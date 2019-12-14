package main

import "crypto/md5"

import "io"

import "fmt"

func main() {
	hasher := md5.New()

	io.WriteString(hasher, "hello")
	io.WriteString(hasher, "world")

	hash := hasher.Sum([]byte("0"))

	fmt.Printf("hash : %x\n", hash)

	byteArray := md5.Sum([]byte("helloworld"))

	fmt.Printf("hash : %x\n", byteArray)

}