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
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Phones    []string  `json:"phones,omitempty"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedBy string    `json:"updatedBy"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Active    bool      `json:"active"`
}
