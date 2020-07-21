// session core model

package session

import "time"

// Session model
type Session struct {
	ID        string //session ID - uuidV4
	UID       string //user id - uuidV4
	CreatedAt time.Time
	ExpiresAt time.Time
}
