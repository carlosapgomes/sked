package services

import (
	"encoding/base64"
	"errors"
	"fmt"
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
		return nil, errors.New("New appointment creation: returned repository ID not equal to new appointment ID")
	}
	return id, err
}

// Update - updates an appointment
func (s *appointmentService) Update(appointmt appointment.Appointment) (*string, error) {
	// get original appointment
	original, err := s.repo.FindByID(appointmt.ID)
	if err != nil {
		fmt.Print("could not find appointment\n")
		return nil, appointment.ErrNoRecord
	}
	_, err = uuid.FromString(appointmt.UpdatedBy)
	if err != nil {
		return nil, appointment.ErrInvalidInputSyntax
	}
	original.DateTime = appointmt.DateTime
	original.Canceled = appointmt.Canceled
	original.Completed = appointmt.Completed
	original.UpdatedAt = time.Now().UTC()
	original.UpdatedBy = appointmt.UpdatedBy
	original.Notes = appointmt.Notes

	id, err := s.repo.Update(*original)
	if err != nil {
		return nil, err
	}
	if (id != nil) && (*id != original.ID) {
		return nil, errors.New("Appointment update: returned repository ID not equal to new appointment ID")
	}
	return id, nil
}

// FindByID - look for an appointment by its uuid
func (s *appointmentService) FindByID(id string) (*appointment.Appointment, error) {
	_, err := uuid.FromString(id)
	if err != nil {
		return nil, appointment.ErrInvalidInputSyntax
	}
	appointment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

// FindFindByPatientID - look for appointments by its patientID
func (s *appointmentService) FindByPatientID(patientID string) ([]*appointment.Appointment, error) {
	_, err := uuid.FromString(patientID)
	if err != nil {
		return nil, appointment.ErrInvalidInputSyntax
	}
	appointmts, err := s.repo.FindByPatientID(patientID)
	if err != nil {
		return nil, err
	}
	return appointmts, nil
}

// FindByDoctorID - look for appointments by doctorID
func (s *appointmentService) FindByDoctorID(doctorID string) ([]*appointment.Appointment, error) {
	_, err := uuid.FromString(doctorID)
	if err != nil {
		return nil, appointment.ErrInvalidInputSyntax
	}
	appointmts, err := s.repo.FindByDoctorID(doctorID)
	if err != nil {
		return nil, err
	}
	return appointmts, nil
}

// FindByDate - look for appointments by date
func (s *appointmentService) FindByDate(dateTime time.Time) ([]*appointment.Appointment, error) {
	appointmts, err := s.repo.FindByDate(dateTime)
	if err != nil {
		return nil, err
	}
	return appointmts, nil
}

// GetAll - return all appointments
func (s *appointmentService) GetAll(before string, after string, pgSize int) (*appointment.Cursor, error) {
	var appointmtsResp appointment.Cursor
	var err error
	var aList *[]appointment.Appointment
	if pgSize <= 0 {
		pgSize = 15
	}

	switch {
	case (before != "" && after != ""):
		// if both (before & after) are present, returns error
		return nil, appointment.ErrInvalidInputSyntax
	case (before == "" && after == ""):
		// if they are empty
		// get default list and page size
		aList, appointmtsResp.HasBefore, err = s.repo.
			GetAll("", false, pgSize)
		if err != nil {
			return nil, err
		}
		if aList != nil {
			for _, a := range *aList {
				appointmtsResp.Appointments = append(appointmtsResp.Appointments, a)
			}
		}
		if len(appointmtsResp.Appointments) > 0 {
			appointmtsResp.Before = base64.StdEncoding.
				EncodeToString([]byte(appointmtsResp.Appointments[len(appointmtsResp.Appointments)-1].ID))
		} else {
			appointmtsResp.Before = ""
		}
		appointmtsResp.After = ""
		appointmtsResp.HasAfter = false
		// and return values
	case (before != ""):
		// if before is present,
		// get a before list
		c, err := base64.StdEncoding.DecodeString(before)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		aList, appointmtsResp.HasBefore, err = s.repo.GetAll(cursor, false, pgSize)
		if err != nil {
			return nil, err
		}
		if aList != nil {
			for _, a := range *aList {
				appointmtsResp.Appointments = append(appointmtsResp.Appointments, a)
			}
		}
		if len(appointmtsResp.Appointments) > 0 {
			befCursor := base64.StdEncoding.EncodeToString([]byte(appointmtsResp.Appointments[len(appointmtsResp.Appointments)-1].ID))
			appointmtsResp.Before = befCursor
		} else {
			appointmtsResp.Before = ""
		}
		// test for 'after data' from the requested cursor
		// fill the response fields
		_, appointmtsResp.HasAfter, err = s.repo.GetAll(cursor, true, pgSize)
		if appointmtsResp.HasAfter {
			appointmtsResp.After = base64.StdEncoding.EncodeToString([]byte(before))
		} else {
			appointmtsResp.After = ""
		}
		// and return it
	case (after != ""):
		// if after is present,
		// get an after list
		c, err := base64.StdEncoding.DecodeString(after)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		aList, appointmtsResp.HasAfter, err = s.repo.
			GetAll(cursor, true, pgSize)
		// and return it
		if aList != nil {
			for _, a := range *aList {
				appointmtsResp.Appointments = append(appointmtsResp.Appointments, a)
			}
		}
		if len(appointmtsResp.Appointments) > 0 {
			appointmtsResp.After = base64.StdEncoding.EncodeToString([]byte(appointmtsResp.Appointments[0].ID))
		}
		// test for 'before data' from the requested cursor
		// fill the response fields
		_, appointmtsResp.HasBefore, err = s.repo.
			GetAll(cursor, false, pgSize)
		if appointmtsResp.HasBefore {
			appointmtsResp.Before = base64.StdEncoding.EncodeToString([]byte(after))
		}
	}
	return &appointmtsResp, nil
}
