package main

import (
	"RyanFin/netAPI/pkg/model"
	"encoding/json"
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
	http.HandleFunc("/users/{user_id}", getUser)
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
	// In the Go net/http package, you cannot directly set the HTTP method type for http.HandleFunc.
	// However, you can achieve this by using http.HandlerFunc instead and checking the request method within your handler.
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data := loadDataFromFile("users.json")

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

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// logic to retrieve a specific user
	id := r.PathValue("user_id")

	data := loadDataFromFile("users.json")

	var users []model.User

	err := json.Unmarshal(data, &users)
	if err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println(users)

	var user model.User

	for _, e := range users {
		// search for record
		if id == e.ID {
			user = e
		}
	}

	if user.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// I need to write the object directly into a JSON response
	encoder := json.NewEncoder(w)

	// the call to WriteHeader(http.StatusOK) is unnecessary when using json.Encoder
	// because it automatically writes the HTTP status code of 200 (OK) for you.
	err = encoder.Encode(user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func loadDataFromFile(fileName string) []byte {

	userFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer userFile.Close()

	data, err := io.ReadAll(userFile)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data

}
