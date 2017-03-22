package command

import (
	"encoding/json"
	"net/http"

	"github.com/CactusDev/CactusAPI-Go/rethink"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/Google/uuid"
	logger "github.com/Sirupsen/logrus"
)

// PatchHandler handles all PATCH requests to /:token/command/:commandName
func PatchHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	res, _ := json.Marshal(vars)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Handler handles all requests to list commands
func Handler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	conn := rethink.Connection{
		DB:    "cactus",
		Table: "commands",
		URL:   "localhost",
	}
	err := conn.Connect()
	if err != nil {
		logger.Error(err.Error())
		logger.Error("Failed to connect to RethinkDB! Is it running?")
		http.Error(w, "Internal Server Error!", 500)
		return
	}
	response, err := conn.GetTable()
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error!", 500)
		return
	}

	m := schemas.Message{
		Data: schemas.Data{
			ID:         uuid.New().String(),
			Attributes: response,
			Type:       "command",
		},
	}
	res, err := json.Marshal(m)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
