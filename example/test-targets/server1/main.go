package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(1 * time.Second)
		io.WriteString(w, req.Host+"\n")
	}
	port := ":8082"
	http.HandleFunc("/", helloHandler)
	fmt.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
