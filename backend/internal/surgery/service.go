package surgery

import "time"

// Service interface for surgery  model
type Service interface {
	Create(dateTime time.Time,
		PatientName string,
		PatientID string,
		DoctorName string,
		DoctorID string,
		Notes string,
		ProposedSurgery, createdBy string) (*string, error)
	Update(surgery Surgery) (*string, error)
	FindByID(id string) (*Surgery, error)
	FindByPatientID(patientID string) ([]*Surgery, error)
	FindByDoctorID(doctorID string) ([]*Surgery, error)
	FindByDate(date time.Time) ([]*Surgery, error)
	GetAll(before string, after string, pgSize int) (*Cursor, error)
}
