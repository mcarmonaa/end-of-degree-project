package message

// Encrypted is the type message sent and received. It wraps any other message in the payload.
type Encrypted struct {
	Mail    string `json:"mail,omitempty"`
	IVector string `json:"iv"`
	Payload string `json:"payload"`
}
