package main

import (
	"fmt"
	"encoding/base64"
)

func main() {
	fmt.Printf("Standard base64 encoding.\n")

	encodedStr := base64.StdEncoding.EncodeToString([]byte("test"))
	fmt.Printf("Standard : %s\n", encodedStr)

	fmt.Printf("URL base64 encoding.\n")
	encodedStr = base64.URLEncoding.EncodeToString([]byte("test"))
	fmt.Printf("URL : %s\n", encodedStr)

}