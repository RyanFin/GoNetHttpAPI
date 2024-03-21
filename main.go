package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

const (
	port = ":8080"
)

func main() {
	http.HandleFunc("/", handleRoot)
	fmt.Printf("listening on port: %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

}
