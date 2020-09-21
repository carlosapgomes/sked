// user core model

package patient

import "time"

// Cursor encapsulates data for pagination
//type Cursor struct {
//Before    string `json:"before"`
//HasBefore bool   `json:"hasbefore"`
//After     string `json:"after"`
//HasAfter  bool   `json:"hasafter"`
//Patients  []Patient
//}

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
	Phones    []string  `json:"phones,omitempty"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedBy string    `json:"updatedBy"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
