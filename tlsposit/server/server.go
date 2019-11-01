// server.go
//
//      prereq: create a server private key and X.509 certificate - place them in a subdirectory called cert
// 		openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 3650 -out server.crt -subj"/C=US/ST=Kansas/L=Lawrence/O=Enterprise Security/OU=CVAS/CN=*"
//      take server.crt file and move it into the cert subdirectory under the tlsclient binary
//
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	var httpServer = http.Server{
		Addr: ":9191",
	}
	var http2Server = http2.Server{}
	_ = http2.ConfigureServer(&httpServer, &http2Server)
	http.HandleFunc("/hello/sayHello", echoPayload)
	log.Printf("\n[*] Go Backend: { HTTPVersion = 2 }; \n[*] serving on https://localhost:9191/hello/sayHello")
	log.Fatal(httpServer.ListenAndServeTLS("./cert/server.crt", "./cert/server.key"))
}

func echoPayload(w http.ResponseWriter, req *http.Request) {
	//log.Printf("Request connection: %s, path: %s", req.Proto, req.URL.Path[1:])
	defer req.Body.Close()
	contents, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Oops! Failed reading body of the request.\n %s", err)
		http.Error(w, err.Error(), 500)
	}
	log.Printf("[*] %s\n", string(contents))
	fmt.Fprintf(w, "%s\n", string(contents))
}
