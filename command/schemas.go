package command

import "time"
import "github.com/asaskevich/govalidator"

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
	Message []MessageSchema `json:"message"`
	User    string          `json:"user"`
	Action  bool            `json:"action"`
	Target  interface{}     `json:"target"`
	Role    int             `json:"role"`
}

// MessageSchema defines the exact schema of the section of the response that is displayed to the user
type MessageSchema struct {
	Type string `json:"type" valid:"string,required"`
	Data string `json:"data" valid:"string,required"`
	Text string `json:"text" valid:"string,required"`
}

govalidator.CustomTypeTagMap.Set("CmdResponse", CustomTypeValidator(func(i interface{}, context interface{}) bool {
  switch v := context.(type) { // you can type switch on the context interface being validated
  case StructWithCustomByteArray:
    // you can check and validate against some other field in the context,
    // return early or not validate against the context at all – your choice
  case SomeOtherType:
    // ...
  default:
    // expecting some other type? Throw/panic here or continue
  }

  switch v := i.(type) { // type switch on the struct field being validated
  case CustomByteArray:
    for _, e := range v { // this validator checks that the byte array is not empty, i.e. not all zeroes
      if e != 0 {
        return true
      }
    }
  }
  return false
}))
// CreationSchema defines the data needed from the user for validaton during command creation
type CreationSchema struct {
	Response Response `json:"response" valid:"CmdResponse, required"`
}

// Validate validates 
func (obj CreationSchema) Validate() error {
	return nil
}
