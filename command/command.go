package command

import (
	"encoding/json"
	"net/http"
	"time"

	"log"

	"fmt"

	"github.com/CactusDev/CactusAPI-Go/rethink"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/Google/uuid"
	"github.com/gorilla/mux"
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

// Handler handles all requests to list commands
func Handler(w http.ResponseWriter, req *http.Request) {
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
	bytes, _ := json.Marshal(response)
	fmt.Println(string(bytes))

	response, err = conn.GetByUUID("b9a960fd-9b31-47e0-84c4-bd05c78e793c")
	if err != nil {
		log.Fatal(err.Error())
	}
	bytes, _ = json.Marshal(response)
	fmt.Println(string(bytes))

	rVars := mux.Vars(req)

	attr := attributes{
		Name: "foo",
		Response: Response{
			User:   "paradigmshift3d",
			Action: false,
			Role:   0,
			Target: nil,
			Message: []MessageSchema{
				MessageSchema{
					Type: "text",
					Data: "foobar test123",
					Text: "foobar test123",
				},
			},
		},
		CreatedAt: time.Now(),
		Token:     rVars["token"],
		Enabled:   true,
		Count:     0,
		Arguments: []MessageSchema{},
	}

	m := schemas.Message{
		Data: schemas.Data{
			ID:         uuid.New().String(),
			Attributes: attr,
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
