package message

// AuthUser is the message sent by a client to register a new user, or request for its salt. Server replies with
// this message for request salt too. Salt field is plained as a base64.
// {
// 	mail: "manuel",
// 	password: "pass"
//	salt: "JAo="
// }
type AuthUser struct {
	Mail     string `json:"mail,omitempty"`
	Password string `json:"password,omitempty"`
	Salt     string `json:"salt,omitempty"`
}

// Encrypted is the message changed in a communication between a logged user and the system. It wraps any other message in the payload.
// The mail field will be omitted in some cases. iv and payload fields will be plained as base64.
// {
// 	mail: "any@mail.cool",
// 	iv: "aW5pdCB2ZWN0b3IK",
// 	payload: "ZW5jcnlwdGVkIHN1cGVyIHNlY3JldCBkYXRhCg=="
// }
type Encrypted struct {
	Mail    string `json:"mail,omitempty"`
	IVector string `json:"iv"`
	Payload string `json:"payload"`
}

// Widespread is the accommodation for any message.
type Widespread struct {
	Timestamp int64  `json:"timestamp"`
	Nonce     uint64 `json:"nonce"`
	Content   string `json:"content"`
}

// Login is the message sent from a client to be logged in the system. It will send as a message.Encrypte's payload ciphered with the
// user's kdf(password, salt)
type Login struct {
	SharedKey string `json:"ks"`
}
