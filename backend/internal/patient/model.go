// user core model

package patient

import "time"

// Page encapsulates data and pagination cursors
type Page struct {
	StartCursor     string    `json:"startCursor"`
	HasPreviousPage bool      `json:"hasPreviousPage"`
	EndCursor       string    `json:"endCursor"`
	HasNextPage     bool      `json:"hasNextPage"`
	Patients        []Patient `json:"patients"`
}

//User type
type Patient struct {
	ID        string    `json:"id"` //uuidv4
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Phones    []string  `json:"phones"`
	CreatedBy string    `json:"createdBy,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedBy string    `json:"updatedBy,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
