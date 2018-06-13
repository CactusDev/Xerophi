package main

import (
	// Default lib imports
	"flag"
	"fmt"
	"net/http"
	"time"

	// Xerophi modules
	"github.com/CactusDev/Xerophi/redis"
	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/secure"
	"github.com/CactusDev/Xerophi/types"

	// Endpoints
	"github.com/CactusDev/Xerophi/command"
	"github.com/CactusDev/Xerophi/quote"
	"github.com/CactusDev/Xerophi/repeat"
	"github.com/CactusDev/Xerophi/social"
	"github.com/CactusDev/Xerophi/trust"

	// Gin imports

	"github.com/gin-gonic/gin"

	// Debug imports
	debugFname "github.com/onrik/logrus/filename"
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
		log.AddHook(debugFname.NewHook())
	}

	// Load the config
	config = LoadConfig()

	// Initialize connection to RethinkDB
	log.Info("Connecting to RethinkDB...")
	rethink.RethinkConn = &rethink.Connection{
		DB:   config.Rethink.DB,
		Opts: config.Rethink.Connection,
	}
	// Validate connection
	if err := rethink.RethinkConn.Connect(); err != nil {
		log.Fatal("RethinkDB Connection Failed! - ", err)
	}
	log.Info("Success!")

	// Initialize connection Redis
	log.Info("Connecting to Redis...")
	redis.RedisConn = &redis.Connection{
		DB:   config.Redis.DB,
		Opts: config.Redis.Connection,
	}
	// Validate connection
	if err := redis.RedisConn.Connect(); err != nil {
		log.Fatal("Redis Connection Failed! - ", err)
	}
	log.Info("Success!")
}

func generateRoutes(resources map[string]types.Handler, api *gin.RouterGroup) {
	for baseName, handler := range resources {
		// Resources that are generally accessible without auth
		open := api.Group(baseName)

		for _, route := range handler.Routes() {
			var currentGroup = open

			if !route.Enabled {
				// Route currently disabled
				continue
			}

			if len(route.Scopes) > 0 {
				// Figure out which scopes are required
				// The group needs to be a separate group that matches this
				// route's scopes
				currentGroup = api.Group(baseName)
				currentGroup.Use(secure.AuthMiddleware(route.Scopes))
			}

			switch route.Verb {
			case "GET":
				currentGroup.GET(route.Path, route.Handler)
			case "PATCH":
				currentGroup.PATCH(route.Path, route.Handler)
			case "POST":
				currentGroup.POST(route.Path, route.Handler)
			case "DELETE":
				currentGroup.DELETE(route.Path, route.Handler)
			}
		}
	}
}

func main() {
	// Initialize the resources with their associated paths
	resources := map[string]types.Handler{
		"/user/:token/command": &command.Command{
			Table: "commands", Conn: rethink.RethinkConn},
		"/user/:token/quote": &quote.Quote{
			Table: "quotes", Conn: rethink.RethinkConn},
		"/user/:token/social": &social.Social{
			Table: "socials", Conn: rethink.RethinkConn},
		"/user/:token/trust": &trust.Trust{
			Table: "trusts", Conn: rethink.RethinkConn},
		"/user/:token/repeat": &repeat.Repeat{
			Table: "repeats", Conn: rethink.RethinkConn},
	}

	// Initialize the router
	router := gin.Default()
	api := router.Group("/api/v2")

	// Initialize panic recovery middleware
	router.Use(gin.Recovery())

	// Intialize the monitoring/status system
	monitor := rethink.Status{
		Tables: map[string]struct{}{
			"commands": {},
			"quotes":   {},
			"socials":  {},
			"trusts":   {},
			"repeats":  {},
		},
		DBs: map[string]struct{}{
			"cactus": {},
		},
		LastUpdated: time.Now(),
	}

	// Initialize JWT authentication
	secure.SetSecret(config.Secure.Secret)

	// TODO: Add a monitor for Redis
	monitor.Monitor(rethink.RethinkConn)
	api.GET("/status", monitor.APIStatusHandler)
	api.POST("/user/:token/login", secure.Authenticator)

	// Load the routes for the individual handlers
	generateRoutes(resources, api)

	// Start up the Gin router on the configured port
	router.Run(fmt.Sprintf(":%d", config.Server.Port))

	// Start the HTTP server
	log.Warnf("API starting on :%d - %s", port, router.BasePath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
