package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(10 * time.Second)
		io.WriteString(w, "Server 1!\n")
	}

	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
