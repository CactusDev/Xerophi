package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
)

var logger = log.New()
var port int

func init() {
	var debug, verbose bool
	flag.BoolVar(&debug, "debug", false, "Run the API in debug mode")
	flag.BoolVar(&debug, "d", false, "Run the API in debug mode")
	flag.BoolVar(&verbose, "verbose", false, "Run the API in verbose mode")
	flag.BoolVar(&verbose, "v", false, "Run the API in verbose mode")
	flag.IntVar(&port, "port", 8000, "Specify which port the API will run on")
	flag.Parse()

	if debug {
		logger.Warn("Starting API in debug mode!")
	} else if verbose {
		logger.Warn("Starting API in verbose mode!")
	}

	if debug || verbose {
		logger.Level = log.DebugLevel
	}

}

func main() {

}
