package mw

type TokenPayloadData struct {
	UserID    string `json:"user_id,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
}
