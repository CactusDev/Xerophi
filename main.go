package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/CactusDev/Xerophi/command"
	"github.com/CactusDev/Xerophi/quote"
	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/types"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

var port int
var config Config

func init() {
	var debug, verbose bool
	flag.BoolVar(&debug, "debug", false, "Run the API in debug mode")
	flag.BoolVar(&debug, "d", false, "Run the API in debug mode")
	flag.BoolVar(&verbose, "verbose", false, "Run the API in verbose mode")
	flag.BoolVar(&verbose, "v", false, "Run the API in verbose mode")
	flag.IntVar(&port, "port", 8000, "Specify which port the API will run on")
	flag.Parse()

	if debug {
		log.Warn("Starting API in debug mode!")
		gin.SetMode(gin.DebugMode)
	} else if verbose {
		log.Warn("Starting API in verbose mode!")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if debug || verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Load the config
	config = LoadConfig()
}

func generateRoutes(h types.Handler, g *gin.RouterGroup) {
	for _, r := range h.Routes() {
		if !r.Enabled {
			// Route currently disabled
			continue
		}
		switch r.Verb {
		case "GET":
			g.GET(r.Path, r.Handler)
		case "PATCH":
			g.PATCH(r.Path, r.Handler)
		case "POST":
			g.POST(r.Path, r.Handler)
		case "DELETE":
			g.DELETE(r.Path, r.Handler)
		}
	}
}

func main() {
	rdbConn := rethink.Connection{
		DB:   config.Rethink.DB,
		Opts: config.Rethink.Connection,
	}
	err := rdbConn.Connect()
	if err != nil {
		log.Fatal("RethinkDB Connection Failed! - ", err)
	}

	handlers := map[string]types.Handler{
		"/user/:token/command": &command.Command{
			Conn:  &rdbConn,
			Table: "commands",
		},
		"/user/:token/quote": &quote.Quote{
			Conn:  &rdbConn,
			Table: "quotes",
		},
	}

	router := gin.Default()
	api := router.Group("/api/v1")

	// Intialize the monitoring/status system
	monitor := rethink.Status{
		Tables: map[string]struct{}{
			"commands": {},
		},
		DBs: map[string]struct{}{
			"cactus": {},
		},
		LastUpdated: time.Now(),
	}

	monitor.Monitor(&rdbConn)
	api.GET("/status", monitor.APIStatusHandler)

	for baseRoute, handler := range handlers {
		group := api.Group(baseRoute)
		generateRoutes(handler, group)
	}

	router.Run(fmt.Sprintf(":%d", config.Server.Port))

	log.Warnf("API starting on :%d - %s", port, router.BasePath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
