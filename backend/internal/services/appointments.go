package services

import (
	"errors"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
	uuid "github.com/satori/go.uuid"
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
func (s *appointmentService) Create(dateTime time.Time, patientName, patientID, doctorName, doctorID, notes, createdBy string) (*string, error) {
	// validate ID format (uuidV4)
	ptID, err := uuid.FromString(patientID)
	if err != nil {
		return nil, appointment.ErrInvalidInputSyntax
	}
	docID, err := uuid.FromString(doctorID)
	if err != nil {
		return nil, appointment.ErrInvalidInputSyntax
	}
	createdByID, err := uuid.FromString(createdBy)
	if err != nil {
		return nil, appointment.ErrInvalidInputSyntax
	}
	uid := uuid.NewV4()
	dt := dateTime.UTC()
	newAppointmt := appointment.Appointment{
		ID:          uid.String(),
		DateTime:    dt,
		PatientName: patientName,
		PatientID:   ptID.String(),
		DoctorName:  doctorName,
		DoctorID:    docID.String(),
		Notes:       notes,
		Canceled:    false,
		Completed:   false,
		CreatedBy:   createdByID.String(),
		CreatedAt:   time.Now().UTC(),
	}

	id, err := s.repo.Create(newAppointmt)
	if err != nil {
		return nil, err
	}
	if (id != nil) && (*id != newAppointmt.ID) {
		return nil, errors.New("New appointment creation: returned repository ID not equal to new user ID")
	}
	return id, err
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
	return appoints, nil
}

// FindByDoctorID - look for appointments by doctorID
func (s *appointmentService) FindByDoctorID(doctorID string) ([]*appointment.Appointment, error) {
	var appoints []*appointment.Appointment
	return appoints, nil
}

// FFindByDate - look for appointments by date
func (s *appointmentService) FindByDate(date time.Time) ([]*appointment.Appointment, error) {
	var appoints []*appointment.Appointment
	return appoints, nil
}

// GetAll - return all appointments
func (s *appointmentService) GetAll(before string, after string, pgSize int) (*appointment.Cursor, error) {
	var cursor appointment.Cursor
	return &cursor, nil
}
