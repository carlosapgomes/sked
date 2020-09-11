// user core model

package user

import "time"

const (
	// RoleAdmin role const
	RoleAdmin string = "Admin"
	// RoleCommon role const
	RoleCommon string = "Common"
	// RoleDoctor role const
	RoleDoctor string = "Doctor"
)

// Cursor encapsulates data for pagination
type Cursor struct {
	Before    string `json:"before"`
	HasBefore bool   `json:"hasbefore"`
	After     string `json:"after"`
	HasAfter  bool   `json:"hasafter"`
	Users     []User
}

//User type
type User struct {
	ID                string    `json:"id"` //uuidv4
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone,omitempty"`
	HashedPw          []byte    `json:"-"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
	Active            bool      `json:"active"`
	EmailWasValidated bool      `json:"emailWasValidated"`
	Roles             []string  `json:"roles,omitempty"`
}
