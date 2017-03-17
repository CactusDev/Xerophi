package schemas

// Message is the base for the API's JSON response
type Message struct {
	Data interface{} `json:"data"`
}
