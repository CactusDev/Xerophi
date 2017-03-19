package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CactusDev/CactusAPI-Go/command"
	"github.com/CactusDev/CactusAPI-Go/quotes"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	mux "github.com/dimfeld/httptreemux"
)

// HomeHandler handles all requests to the base URL
func HomeHandler(w http.ResponseWriter, req *http.Request, _ map[string]string) {
	m := schemas.Message{
		Data: "Ohai! You're home!",
	}
	response, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// HomeOptions returns an array of the HTTP request options available for this endpoint
func HomeOptions(w http.ResponseWriter, req *http.Request, _ map[string]string) {
	m := [2]string{"GET", "OPTIONS"}
	response, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	router := mux.New()
	api := router.NewGroup("/api/v1")

	router.GET("/", HomeHandler)
	router.OPTIONS("/", HomeOptions)
	api.GET("/:token/command", command.Handler)
	api.PATCH("/:token/command/:commandName", command.PatchHandler)
	api.GET("/:token/quote", quotes.Handler)
	log.Fatal(http.ListenAndServe(":8000", router))
}
