package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CactusDev/CactusAPI-Go/command"
	"github.com/CactusDev/CactusAPI-Go/quotes"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/Google/uuid"
	"github.com/gorilla/mux"
)

// HomeHandler handles all requests to the base URL
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	m := schemas.Message{
		Data: schemas.Data{
			ID:         uuid.New().String(),
			Attributes: "stuff",
			Type:       "nil",
		},
	}
	response, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/{token}/command", command.Handler).Methods("GET", "OPTIONS")
	r.HandleFunc("/{token}/quote", quotes.Handler).Methods("GET", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8000", r))
}
