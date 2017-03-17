package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Google/uuid"
	"github.com/gorilla/mux"
)

// Message is the base for the API's JSON response
type Message struct {
	Data string `json:"data"`
	ID   string `json:"id,omitempty"`
}

// HomeHandler handles all requests to the base URL
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	m := Message{
		Data: "Ohai, you're home!",
		ID:   uuid.New().String(),
	}
	response, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// CommandHandler handles all requests to list commands
func CommandHandler(w http.ResponseWriter, req *http.Request) {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
