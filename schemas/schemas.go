package schemas

// MessagePacket is the low-level format for storing contents of a message and arguments
type MessagePacket struct {
	Data string `json:"data" validate:"required"`
	Text string `json:"text" validate:"required"`
	Type string `json:"type" validate:"required"`
}
