package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":12970", nil)
}

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
