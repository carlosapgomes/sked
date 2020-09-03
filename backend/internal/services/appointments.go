package services

import (
	"time"

	"carlosapgomes.com/sked/internal/appointment"
)

// AppointmentService provides implementation of appointment domain interface
type appointmentService struct {
	repo appointment.Repository
}

// NewAppointmentService returns an appointment Service instance
func NewAppointmentService(repo appointment.Repository) appointment.Service {
	return &appointmentService{
		repo,
	}
}

// Create - creates a new appointment and returns its uuid
func (s *appointmentService) Create(dateTime time.Time, patientName, patientID, doctorName, doctorID, notes string) (*string, error) {
	var id string
	return &id, nil
}

// Update - updates an appointment
func (s *appointmentService) Update(appointment appointment.Appointment) (*string, error) {
	var id string
	return &id, nil
}

// FindByID - look for an appointment by its uuid
func (s *appointmentService) FindByID(id string) (*appointment.Appointment, error) {
	var appointment appointment.Appointment
	return &appointment, nil
}

// FindFindByPatientID - look for appointments by its patientID
func (s *appointmentService) FindByPatientID(patientID string) ([]*appointment.Appointment, error) {
	var appoints []*appointment.Appointment
	return &appoints, nil
}

// FindFindByDoctorID - look for appointments by doctorID
func (s *appointmentService) FindByDoctorID(doctorID string) ([]*appointments.Appointment, error) {
	var appoints []*appointment.Appointment
	return &appoints, nil
}

// FFindByDate - look for appointments by date
func (s *appointmentService) FindByDate(date time.Time) ([]*appointment.Appointment, error) {
	var appoints []*appointment.Appointment
	return &id, nil
}

// GetAll - return all appointments
func (s *appointmentService) GetAll(before string, after string, pgSize int) (*appointment.Cursor, error) {
	var cursor appointment.Cursor
	return &cursor, nil
}
