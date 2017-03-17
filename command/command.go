package command

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/Google/uuid"
)

// Command foo
type Command struct {
	ID         string     `json:"id"`
	Attributes attributes `json:"attributes"`
	Type       string     `json:"type"`
}

// Attributes is a schema for JSON responses regarding commands
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
	// rVars := mux.Vars(req)

	message := []MessageSchema{
		MessageSchema{
			Type: "text",
			Data: "foobar test123",
			Text: "foobar test123",
		},
	}

	response := Response{
		User:    "paradigmshift3d",
		Action:  false,
		Role:    0,
		Target:  nil,
		Message: message,
	}

	attr := attributes{
		Name:      "foo",
		Response:  response,
		CreatedAt: time.Now(),
		Token:     "paradigmshift3d",
		Enabled:   true,
		Count:     0,
		Arguments: []MessageSchema{},
	}

	m := schemas.Message{
		Data: Command{
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
