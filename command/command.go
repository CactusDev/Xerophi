package command

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/CactusDev/CactusAPI-Go/rethink"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/Google/uuid"
)

// PatchHandler handles all PATCH requests to /:token/command/:commandName
func PatchHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	res, _ := json.Marshal(vars)
	w.Write(res)
	w.WriteHeader(http.StatusAccepted)
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
		log.Fatal(err.Error())
	}
	response, err := conn.GetTable()
	if err != nil {
		log.Fatal(err.Error())
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
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
