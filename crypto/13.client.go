package main

import (
 "io/ioutil"
 "log"
 "crypto/x509"
 "crypto/tls"
 "fmt"
 "net/http"
)

func main() {
	caCertInfo, err := ioutil.ReadFile("./server.crt")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertInfo)

	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")

	tlsCfg := tls.Config {
		RootCAs: certPool,
		Certificates: []tls.Certificate{clientCert},
	}

	client := http.Client {
		Transport: &http.Transport {
			TLSClientConfig: &tlsCfg,
		},
	}

	response, err := client.Get("https://localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	bodyInfo, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	fmt.Printf("body : %s\n", bodyInfo)
	fmt.Printf("status code : %s\n", response.Status)

}