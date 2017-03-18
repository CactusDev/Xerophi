package quotes

import (
	"net/http"

	"encoding/json"

	"github.com/CactusDev/CactusAPI-Go/schemas"
	"github.com/Google/uuid"
)

type attributes struct {
	Quote   string `json:"quote"`
	QuoteID int    `json:"quoteId"`
	Token   string `json:"token"`
}

// Handler handles all requests to the /{token}/quote endpoint
func Handler(w http.ResponseWriter, req *http.Request, rVars map[string]string) {
	attr := attributes{
		Quote:   "Spam eggs!",
		QuoteID: 1,
		Token:   rVars["token"],
	}

	m := schemas.Message{
		Data: schemas.Data{
			ID:         uuid.New().String(),
			Attributes: attr,
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
