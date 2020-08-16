// user core model

package patient

import "time"

// Cursor encapsulates data for pagination
type Cursor struct {
	Before    string `json:"before"`
	HasBefore bool   `json:"hasbefore"`
	After     string `json:"after"`
	HasAfter  bool   `json:"hasafter"`
	Users     []Patient
}

//User type
type Patient struct {
	ID        string    `json:"id"` //uuidv4
	Name      string    `json:"name"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Active    bool      `json:"active"`
}
