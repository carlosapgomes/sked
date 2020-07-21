package token

import "time"

const (
	// PwReset token type constant
	PwReset string = "PwReset"
	// ValidateEmail token type constant
	ValidateEmail string = "ValidateEmail"

	// ExpPwReset PwReset token expiration time in seconds
	ExpPwReset time.Duration = 30 * time.Minute // 30 min
	// ExpValidateEmail ValidateEmail token expiration time in seconds
	ExpValidateEmail time.Duration = 1 * time.Hour // 01 hour
)

//Token model
type Token struct {
	ID        string    // random token
	UID       string    //user id string (uuidV4)
	CreatedAt time.Time // UTC
	ExpiresAt time.Time // UTC
	Kind      string
}
