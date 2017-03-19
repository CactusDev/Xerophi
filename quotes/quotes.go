package quotes

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/CactusDev/CactusAPI-Go/rethink"
	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/CactusDev/CactusAPI-Go/util"
)

var conn = rethink.Connection{
	DB:    "cactus",
	Table: "quotes",
	URL:   "localhost",
}

// PostHandler handles POST requests to /:token/quote endpoint
func PostHandler(w http.ResponseWriter, req *http.Request, vars map[string]string) {

}

// Handler handles all requests to the /:token/quote endpoint
func Handler(w http.ResponseWriter, req *http.Request, vars map[string]string) {
	err := conn.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	response, err := conn.GetTable()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, quote := range response {
		attr := attributes{}
		err := util.FillStruct(quote.(map[string]interface{}), &attr)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}

	m := schemas.Message{
		Data: schemas.Data{
			Attributes: response,
			Type:       "quote",
		},
	}

	res, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
