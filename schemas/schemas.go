package schemas

// MessagePacket is the low-level format for storing contents of a message and arguments
type MessagePacket struct {
	Data string `jsonapi:"attr,data" json:"data" validate:"required"`
	Text string `jsonapi:"attr,text" json:"text" validate:"required"`
	Type string `jsonapi:"attr,type" json:"type" validate:"required"`
}
