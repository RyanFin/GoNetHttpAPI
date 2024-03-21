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
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	port = ":8080"
)

var (
	// define a global JWT secret key
	jwtKey = []byte("your-secret-key") // Change this to your secret key
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/users", authMiddleware(getAllUsers))
	http.HandleFunc("/users/{user_id}", authMiddleware(getUser))
	// authentication route
	http.HandleFunc("/login", loginHandler)
	fmt.Printf("listening on port: %s...", port)

	s := &http.Server{
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

type auth struct {
	Token string `json:"token"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	// Here you should implement the login logic.
	// For simplicity, let's just assume the user is authenticated and generate a token.
	token := generateToken()
	w.Header().Set("Authorization", "Bearer "+token)

	// fmt.Fprintf(w, "Token generated successfully\n")

	auth := auth{Token: token}
	encoder := json.NewEncoder(w)

	// the call to WriteHeader(http.StatusOK) is unnecessary when using json.Encoder
	// because it automatically writes the HTTP status code of 200 (OK) for you.
	err := encoder.Encode(auth)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// code to generate jwtToken
func generateToken() string {
	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "John Doe"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // token expiration time

	// sign and get the complete encoded token as a string
	tokenString, _ := token.SignedString(jwtKey)

	return tokenString
}

// Middleware function to verify JWT token
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization Header
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure that the token method conforms to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized access: %v", err)
			return
		}

		// Check if token is valid
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized access: invalid token")
			return
		}

		// Proceed to the protected handler if token is valid
		next.ServeHTTP(w, r)
	}
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
