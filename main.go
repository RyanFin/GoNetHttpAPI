package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	port = ":8080"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/users", getAllUsers)
	fmt.Printf("listening on port: %s...", port)

	s := &http.Server{
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userFile, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer userFile.Close()

	data, err := io.ReadAll(userFile)
	if err != nil {
		fmt.Println(err.Error())
	}

	// I don't need to convert the JSON data into a JSON object
	// var users model.Users

	// err = json.Unmarshal(dat, &users)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	w.WriteHeader(http.StatusOK)
	// data must be sent as bytes
	w.Write(data)
}
