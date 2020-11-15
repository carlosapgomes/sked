// appointment repository port interface definition

package appointment

import "time"

// Repository inteface definition for appointment model
type Repository interface {
	Create(appointment Appointment) (*string, error)
	Update(appointment Appointment) (*string, error)
	FindByID(id string) (*Appointment, error)
	FindByPatientID(patientID string) ([]Appointment, error)
	FindByDoctorID(doctorID string) ([]Appointment, error)
	FindByDate(date time.Time) ([]Appointment, error)
	FindByInterval(start, end time.Time) ([]Appointment, error)
	GetAll(cursor string, after bool, pgSize int) ([]Appointment, bool, error)
}
