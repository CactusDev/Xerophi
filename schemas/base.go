package schemas

// Data is the JSON-API spec for returning
type Data struct {
	ID         string      `json:"id,omitempty"`
	Attributes interface{} `json:"attributes"`
	Type       string      `json:"type"`
}

// Message is the base for the API's JSON response
type Message struct {
	Data interface{} `json:"data,omitempty"`
	Meta Meta        `json:"meta,omitempty"`
}

// Meta can take any type as extra information related to the response
type Meta interface{}
