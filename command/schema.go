package command

import (
	"reflect"

	"github.com/CactusDev/Xerophi/schemas"

	log "github.com/Sirupsen/logrus"
)

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID        string                  `jsonapi:"id,primary"`
	Arguments []schemas.MessagePacket `jsonapi:"arguments,attributes"`
	Count     int                     `jsonapi:"count,attributes"`
	CreatedAt string                  `jsonapi:"createdAt,attributes"`
	Enabled   bool                    `jsonapi:"enabled,attributes"`
	Name      string                  `jsonapi:"name,attributes"`
	Response  responseSchema          `jsonapi:"response,attributes"`
	Token     string                  `jsonapi:"token,attributes"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
}

type responseSchema struct{}

// GetAPITag returns the jsonapi tag value for the specified field
func (r ResponseSchema) GetAPITag(lookup string) string {
	field, ok := reflect.TypeOf(r).FieldByName(lookup)
	if !ok {
		log.Warn("Uh, stuff happened")
		return ""
	}
	tag, ok := field.Tag.Lookup("jsonapi")
	if !ok {
		log.Warn("Aw snap")
		return ""
	}
	return tag
}
