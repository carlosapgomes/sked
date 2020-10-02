package storage

import (
	"database/sql"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
)

// appointmentRepository type
type appointmentRepository struct {
	DB *sql.DB
}

// NewPgAppointmentRepository returns an instance of an appointmentRepository
func NewPgAppointmentRepository(db *sql.DB) appointment.Repository {
	return &appointmentRepository{
		db,
	}
}

// Create - creates a new appointment
func (r appointmentRepository) Create(appointmt appointment.Appointment) (*string, error) {
	return nil, nil
}

// Update - updates an appointment
func (r appointmentRepository) Update(appointmt appointment.Appointment) (*string, error) {
	return nil, nil
}

// FindByID - finds an appointment by its ID
func (r appointmentRepository) FindByID(id string) (*appointment.Appointment, error) {
	return nil, nil
}

// FindByPatientID - finds an appointment by its patientID
func (r appointmentRepository) FindByPatientID(patientID string) ([]*appointment.Appointment, error) {
	return nil, nil
}

// FindFindByDoctorID - finds an appointment by its doctorID
func (r appointmentRepository) FindByDoctorID(doctorID string) ([]*appointment.Appointment, error) {
	return nil, nil
}

// FindByDate - finds appointments in a date
func (r appointmentRepository) FindByDate(date time.Time) ([]*appointment.Appointment, error) {
	return nil, nil
}

// GetAll - returns a paginated list of appointments
func (r appointmentRepository) GetAll(cursor string, after bool, pgSize int) (*[]appointment.Appointment, bool, error) {
	return nil, false, nil
}
