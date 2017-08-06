package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/CactusDev/Xerophi/command"
	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/types"

	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
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
	g.GET("", h.GetAll)
	g.PATCH("/:name", h.Update)
	g.GET("/:name", h.GetSingle)
	g.DELETE("/:name", h.Create)
}

func main() {
	rdbConn := rethink.Connection{
		DB:   config.Rethink.DB,
		Opts: config.Rethink.Connection,
	}
	rdbConn.Connect()

	handlers := map[string]types.Handler{
		"/user/:token/command": &command.Command{
			Conn:  &rdbConn,
			Table: "commands",
		},
	}

	router := gin.Default()
	api := router.Group("/api/v1")

	for baseRoute, handler := range handlers {
		group := api.Group(baseRoute)
		generateRoutes(handler, group)
	}

	router.Run(fmt.Sprintf(":%d", config.Server.Port))

	log.Warnf("API starting on :%d - %s", port, router.BasePath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
