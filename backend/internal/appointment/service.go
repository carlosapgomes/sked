package appointment

import "time"

// Service interface for appointment model
type Service interface {
	Create(dateTime time.Time,
		PatientName string,
		PatientID string,
		DoctorName string,
		DoctorID string,
		Notes,
		CreatedBy string) (*string, error)
	Update(appointment Appointment) (*string, error)
	FindByID(id string) (*Appointment, error)
	FindByPatientID(patientID string) ([]Appointment, error)
	FindByDoctorID(doctorID string) ([]Appointment, error)
	FindByDate(date time.Time) ([]Appointment, error)
	GetAll(before string, after string, pgSize int) (*Page, error)
}
