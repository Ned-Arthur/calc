package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"html/template"
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

func main() {
	router := http.NewServeMux()

	// Serve routes for the interface
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/getkey", handleKeyPost)
	
	// Define the routes we'll serve for the API
	router.HandleFunc("POST /add", handleAdd)
	router.HandleFunc("POST /subtract", handleSubtract)
	router.HandleFunc("POST /multiply", handleMultiply)
	router.HandleFunc("POST /divide", handleDivide)


	fmt.Println("Now listening on http://localhost:8080")
	http.ListenAndServe(":8080", router)
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
	}


	tmpl := template.Must(template.ParseFiles("pages/getkey.html"))
	data := APIKeyPageData{
		APIKey: key,
		Exists: exists,
	}
	tmpl.Execute(w, data)
}


func handleAdd(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("400 - Can't divide by 0"))
		return
	}

	// Calculate output and set as HTTP response
	var output Output
	output.Ans = input.A / input.B
	json.NewEncoder(w).Encode(output)
}

