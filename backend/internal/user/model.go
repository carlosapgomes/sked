// user core model

package user

import "time"

const (
	// RoleAdmin role const
	RoleAdmin string = "Admin"
	// RoleClerk role const
	RoleClerk string = "Clerk"
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

// Page encapsulates data and pagination cursors
type Page struct {
	StartCursor     string `json:"startCursor"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	EndCursor       string `json:"endCursor"`
	HasNextPage     bool   `json:"hasNextPage"`
	Users           []User `json:"users"`
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
