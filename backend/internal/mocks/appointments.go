package mocks

import (
	"time"

	"carlosapgomes.com/sked/internal/appointment"
)

// mocks appointment services and repository

// AppointmentMockRepo is a mocked appointment repository
type AppointmentMockRepo struct {
	aDb []appointment.Appointment
}

//NewUserRepo returns a mocked repository
func NewAppointmentRepo() *AppointmentMockRepo {
	var db []appointment.Appointment
	return &AppointmentMockRepo{
		db,
	}
}

// Create
func (r AppointmentMockRepo) Create(appointment appointment.Appointment) (*string, error) {
	var id string
	return &id, nil
}

// Update
func (r AppointmentMockRepo) Update(appointment appointment.Appointment) (*string, error) {
	var id string
	return &id, nil
}

// FindFindByID
func (r AppointmentMockRepo) FindByID(id string) (*appointment.Appointment, error) {
	var appointment appointment.Appointment
	return &appointment, nil
}

// FindByPatientID
func (r AppointmentMockRepo) FindByPatientID(patientID string) ([]*appointment.Appointment, error) {
	var res []*appointment.Appointment
	return res, nil
}

// FindByDoctorID
func (r AppointmentMockRepo) FindByDoctorID(doctorID string) ([]*appointment.Appointment, error) {
	var res []*appointment.Appointment
	return res, nil
}

// FFindByDate
func (r AppointmentMockRepo) FindByDate(date time.Time) ([]*appointment.Appointment, error) {
	var res []*appointment.Appointment
	return res, nil
}

// GetAll
func (r AppointmentMockRepo) GetAll(cursor string, after bool, pgSize int) ([]*appointment.Appointment, bool, error) {
	var res []*appointment.Appointment
	return res, false, nil
}
