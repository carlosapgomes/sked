// surgery repository port interface definition

package surgery

import "time"

// Repository inteface definition for surgery model
type Repository interface {
	Create(surgery Surgery) (*string, error)
	Update(surgery Surgery) (*string, error)
	FindByID(id string) (*Surgery, error)
	FindByPatientID(patientID string) ([]*Surgery, error)
	FindByDoctorID(doctorID string) ([]*Surgery, error)
	FindByDate(date time.Time) ([]*Surgery, error)
	GetAll(cursor string, after bool, pgSize int) (*[]Surgery, bool, error)
}
