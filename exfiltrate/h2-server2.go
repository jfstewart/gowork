//h2-server-ubuntu.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s user-agent=%s", r.Method, r.Proto, r.URL.String(), r.UserAgent())
	if r.Method != http.MethodPost {
			http.Error(w, "only POST method is allowed", http.StatusBadRequest)
			return
	}
	log.Printf("parsing multipart form")
	// no more than 100MB of memory, the rest goes into /tmp
	r.ParseMultipartForm(100000000)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
			http.Error(w, "failed to read uploadfile form field", http.StatusBadRequest)
			log.Printf("failed to read uploadfile form field: %v", err)
			return
	}
	defer file.Close()
	log.Println(handler.Filename)
}

func main() {
	// Create a server on port 12970
	// Exactly how you would run an HTTP1.1 server
	//srv := &http.Server{Addr: ":12970", Handler: http.HandlerFunc(handle)}
	srv := &http.Server{Addr: ":12970"}

	//
	http.HandleFunc("/upload", uploadHandler)
	
	// Start the server with TLS, since we are running HTTP2 it must be run with TLS.
	// Exactly how you would run an HTTP1.1 server with TLS connection.
	log.Printf("Serving on https://0.0.0.0:12970")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))

}
