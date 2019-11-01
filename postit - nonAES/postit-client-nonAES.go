// postit.go - post to a URL
package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
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

func postLoot(lootPtr *string, hostPtr *string, portPtr *string) {

	log.Println("Parsing loot file:", *lootPtr)

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
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(jsonData))
			resp, err := http.Post(s, "application/json", bytes.NewBuffer(jsonData))
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
