package command

import "time"

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
