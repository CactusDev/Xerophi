package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CactusDev/CactusAPI-Go/command"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/gorilla/mux"
)

// HomeHandler handles all requests to the base URL
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	m := schemas.Message{
		Data: "spam",
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
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/command/{commandName}", command.Handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
