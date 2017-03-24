package command

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
)

type dbCommand struct {
	count     int       `gorethink:"count" valid:"int,required"`
	createdAt time.Time `gorethink:"createdAt" valid:"date,required"`
	enabled   bool      `gorethink:"enabled" valid:"bool,required"`
	name      string    `gorethink:"name" valid:"string,required"`
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
	Message []MessageSchema `json:"message" valid:"-"`
	User    string          `json:"user" valid:"string,required"`
	Action  bool            `json:"action" valid:"bool,optional"`
	Target  interface{}     `json:"target" valid:"-,optional"`
	Role    int             `json:"role" valid:"int"`
}

// MessageSchema defines the exact schema of the section of the response that is displayed to the user
type MessageSchema struct {
	Type string `json:"type" valid:"string,required"`
	Data string `json:"data" valid:"string,required"`
	Text string `json:"text" valid:"string,required"`
}

// CreationSchema defines the data needed from the user for validaton during command creation
type CreationSchema struct {
	Response Response `json:"response" valid:"CmdResponse, required"`
}

// Validate validates
func (obj CreationSchema) Validate() (bool, error) {
	return govalidator.ValidateStruct(obj)
}

func init() {
	govalidator.CustomTypeTagMap.Set("CmdResponse", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		fmt.Println(i)
		fmt.Println(o)
		return true
	}))
}
