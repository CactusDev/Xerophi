package main

import "time"

// Command is a schema for JSON responses regarding commands
type Command struct {
	Message   Message
	Name      string    `json:"name"`
	Response  Response  `json:"response"`
	CreatedAt time.Time `json:"createdAt"`
	Token     string    `json:"token"`
	Enabled   bool      `json:"enabled"`
	Arguments []Message `json:"arguments,omitempty"`
	Count     int       `json:"count"`
}

// Response keeps track of the information required for the command's response
type Response struct {
	Message []Message `json:"message"`
	User    string    `json:"user"`
	Action  bool      `json:"action"`
	Target  string    `json:"target"`
	Role    int       `json:"role"`
}

// Message defines the exact schema of the section of the response that is displayed to the user
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Text string `json:"text"`
}
