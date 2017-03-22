package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/CactusDev/CactusAPI-Go/command"
	"github.com/CactusDev/CactusAPI-Go/quotes"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	log "github.com/Sirupsen/logrus"
	mux "github.com/dimfeld/httptreemux"
)

var logger = log.New()
var port = new(int)

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

func init() {
	debug := flag.Bool("debug", false, "Run the API in debug mode")
	verbose := flag.Bool("v", false, "Run the API in verbose mode")
	port = flag.Int("port", 8000, "Specify which port the API will run on")
	flag.Parse()

	if *debug {
		logger.Warn("Starting API in debug mode!")
	} else if *verbose {
		logger.Warn("Starting API in verbose mode!")
	}

	if *debug || *verbose {
		logger.Level = log.DebugLevel
	}
}

func main() {
	router := mux.New()
	api := router.NewGroup("/api/v1")

	logger.WithField("port", *port).Info("Starting API server!")

	router.GET("/", HomeHandler)
	router.OPTIONS("/", HomeOptions)
	api.GET("/user/:token/command", command.Handler)
	api.PATCH("/user/:token/command/:commandName", command.PatchHandler)
	api.GET("/user/:token/quote", quotes.Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}
