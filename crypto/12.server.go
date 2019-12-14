package main

import (
	"fmt"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"log"
)

func main() {
	caCertInfo, err := ioutil.ReadFile("./client.crt")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertInfo)

	tlsCfg := tls.Config {
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool,
	}
	server := http.Server {
		Addr: ":3000",
		Handler: nil,
		TLSConfig: &tlsCfg,
	}

	http.HandleFunc("/", func (writer http.ResponseWriter, request *http.Request)  {
		fmt.Println("HandleFunc call!\n")
		writer.Write([]byte("hello world!"))
	})

	err = server.ListenAndServeTLS("server.crt", "server.key")

	if err != nil {
		log.Fatal(err)
	}
}
