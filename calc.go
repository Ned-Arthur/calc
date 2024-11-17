package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"html/template"

	"github.com/Ned-Arthur/calc/middleware"
)

// JSON data structs

type Input struct {
	A float32 `json:"a"`
	B float32 `json:"b"`
}

type Output struct {
	Ans float32 `json:"ans"`
}

// HTML template structs

type APIKeyPageData struct {
	APIKey string
	Exists bool
}

// In-memory cache of valid API keys and their owners
var keyCache = make(map[string]string)
var emailCache = make(map[string]string)

func main() {
	router := http.NewServeMux()

	// Serve routes for the interface
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/getkey", handleKeyPost)
	
	// Define the routes we'll serve for the API
	router.HandleFunc("POST /add/{api_key}", handleAdd)
	router.HandleFunc("POST /subtract/{api_key}", handleSubtract)
	router.HandleFunc("POST /multiply/{api_key}", handleMultiply)
	router.HandleFunc("POST /divide/{api_key}", handleDivide)


	server := http.Server{
		Addr: ":8080",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Now listening on http://localhost:8080")
	server.ListenAndServe()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/home.html"))
	tmpl.Execute(w, "")		// No data for templating, just render page
}

func handleKeyPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.Form["email"][0]

	var exists bool
	var key string

	if val, ok := keyCache[email]; ok {
		exists = true
		key = val
	} else {
		exists = false
		key = generateKey()
		keyCache[email] = key
		emailCache[key] = email
	}


	tmpl := template.Must(template.ParseFiles("pages/getkey.html"))
	data := APIKeyPageData{
		APIKey: key,
		Exists: exists,
	}
	tmpl.Execute(w, data)
}

// Check that the API key is valid, returning a message to describe why it isn't 
func validateAPIKey(r *http.Request) string {
	key := r.PathValue("api_key")
	if len(key) != KEYLENGTH {
		return "400 - API key is incorrect length"
	}
	if _, ok := emailCache[key]; !ok {
		return "400 - invalid API key"
	}

	return ""
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	errcode := validateAPIKey(r)
	if errcode != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errcode))
		return
	}

	// Get input
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate output and set as HTTP response
	var output Output
	output.Ans = input.A + input.B
	json.NewEncoder(w).Encode(output)
}

func handleSubtract(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	errcode := validateAPIKey(r)
	if errcode != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errcode))
		return
	}

	// Get input
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate output and set as HTTP response
	var output Output
	output.Ans = input.A - input.B
	json.NewEncoder(w).Encode(output)
}

func handleMultiply(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	errcode := validateAPIKey(r)
	if errcode != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errcode))
		return
	}

	// Get input
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate output and set as HTTP response
	var output Output
	output.Ans = input.A * input.B
	json.NewEncoder(w).Encode(output)
}

func handleDivide(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	errcode := validateAPIKey(r)
	if errcode != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errcode))
		return
	}

	// Get input
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Spit an error if we can't divide by 0
	if input.B == 0.0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Can't divide by zero"))
		return
	}

	// Calculate output and set as HTTP response
	var output Output
	output.Ans = input.A / input.B
	json.NewEncoder(w).Encode(output)
}

