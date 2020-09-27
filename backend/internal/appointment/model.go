// appointment core model

package appointment

import "time"

// Cursor encapsulates data for pagination
//type Cursor struct {
//Before       string `json:"before"`
//HasBefore    bool   `json:"hasbefore"`
//After        string `json:"after"`
//HasAfter     bool   `json:"hasafter"`
//Appointments []Appointment
//}

// Page encapsulates data and pagination cursors
type Page struct {
	StartCursor     string        `json:"startCursor"`
	HasPreviousPage bool          `json:"hasPreviousPage"`
	EndCursor       string        `json:"endCursor"`
	HasNextPage     bool          `json:"hasNextPage"`
	Appointments    []Appointment `json:"patients"`
}

//Appointment type
type Appointment struct {
	ID          string    `json:"id"` //uuidv4
	DateTime    time.Time `json:"dateTime"`
	PatientName string    `json:"patientName"`
	PatientID   string    `json:"patientID"`
	DoctorName  string    `json:"doctorName"`
	DoctorID    string    `json:"doctorID"`
	Notes       string    `json:"notes"`
	Canceled    bool      `json:"canceled"`
	Completed   bool      `json:"completed"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedBy   string    `json:"updatedBy"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
