package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]

	protocol := "http"
	server := "localhost:8080"

	requestBuilder := new(bytes.Buffer)
	requestBuilder.Write([]byte(protocol))
	requestBuilder.Write([]byte("://"))
	requestBuilder.Write([]byte(server))
	requestBuilder.Write([]byte("/search?find="))
	for _, arg := range args {
		requestBuilder.Write([]byte(arg))
	}

	requestString := requestBuilder.String()

	resp, err := http.Get(requestString)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Status: %s\n", resp.Status)
	}

	if resp.ContentLength > (1 << 20) {
		log.Fatalf("response too large: %d", resp.ContentLength)
	}

	buffer := make([]byte, 1<<10)
	for {
		n, err := resp.Body.Read(buffer)
		fmt.Print(string(buffer[:n]))
		if err == io.EOF || n == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("\n")
}
