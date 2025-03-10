package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Hello!\n")
}

func main() {
	var port string
	flag.StringVar(&port, "port", "1543", "Port to listen on")
	flag.Parse()

	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
