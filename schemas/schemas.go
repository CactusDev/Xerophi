package schemas

// MessagePacket is the low-level format for storing contents of a message and arguments
type MessagePacket struct {
	Data string `json:"data"`
	Text string `json:"text"`
	Type string `json:"type"`
}
