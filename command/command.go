package command

import (
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/CactusDev/CactusAPI-Go/rethink"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/Google/uuid"
)

type dbCommand struct {
	count     int       `gorethink:"count"`
	createdAt time.Time `gorethink:"createdAt"`
	enabled   bool      `gorethink:"enabled"`
	name      string    `gorethink:"name"`
	response  Response  `gorethink:"response"`
	token     string    `gorethink:"token"`
	id        string    `gorethink:"-"`
}

type attributes struct {
	Name      string          `json:"name"`
	Response  Response        `json:"response"`
	CreatedAt time.Time       `json:"createdAt"`
	Token     string          `json:"token"`
	Enabled   bool            `json:"enabled"`
	Arguments []MessageSchema `json:"arguments,omitempty"`
	Count     int             `json:"count"`
}

// Response keeps track of the information required for the command's response
type Response struct {
	Message []MessageSchema `json:"message"`
	User    string          `json:"user"`
	Action  bool            `json:"action"`
	Target  interface{}     `json:"target"`
	Role    int             `json:"role"`
}

// MessageSchema defines the exact schema of the section of the response that is displayed to the user
type MessageSchema struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Text string `json:"text"`
}

// PatchHandler handles all PATCH requests to /:token/command/:commandName
func PatchHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	res, _ := json.Marshal(vars)
	w.Write(res)
	w.WriteHeader(http.StatusAccepted)
}

// Handler handles all requests to list commands
func Handler(w http.ResponseWriter, req *http.Request, rVars map[string]string) {

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
