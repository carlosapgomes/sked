// surgery core model

package surgery

import (
	"time"
)

// Cursor encapsulates data for pagination
type Cursor struct {
	Before    string `json:"before"`
	HasBefore bool   `json:"hasbefore"`
	After     string `json:"after"`
	HasAfter  bool   `json:"hasafter"`
	Surgeries []Surgery
}
type Surgery struct {
	ID              string    `json:"id"` //uuidv4
	DateTime        time.Time `json:"dateTime"`
	PatientName     string    `json:"patientName"`
	PatientID       string    `json:"patientID"`
	DoctorName      string    `json:"doctorName"`
	DoctorID        string    `json:"doctorID"`
	Notes           string    `json:"notes"`
	ProposedSurgery string    `json:"proposedSurgery"`
	Canceled        bool      `json:"canceled"`
	Done            bool      `json:"done"`
	CreatedBy       string    `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	UpdatedBy       string    `json:"updatedBy"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
}
