package services

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/patient"
	"carlosapgomes.com/sked/internal/user"
	uuid "github.com/satori/go.uuid"
)

// PatientService provides implementation of patient domain interface
type patientService struct {
	repo patient.Repository
}

// NewPatientService returns an patient Service instance
func NewPatientService(repo patient.Repository) patient.Service {
	return &patientService{
		repo,
	}
}

// Create creates a new patient and returns its uuid
func (s *patientService) Create(name, address, city, state string,
	phones []string, createdBy string) (*string, error) {
	if len(name) == 0 {
		return nil, patient.ErrInvalidInputSyntax
	}
	if len(phones) > 10 ||
		len(phones) == 0 {
		return nil, patient.ErrInvalidInputSyntax
	}
	if strings.TrimSpace(strings.Join(phones, "")) == "" {
		return nil, patient.ErrInvalidInputSyntax
	}
	_, err := uuid.FromString(createdBy)
	if err != nil {
		return nil, patient.ErrInvalidInputSyntax
	}
	newID := uuid.NewV4()
	dt := time.Now().UTC()
	newPatient := patient.Patient{
		ID:        newID.String(),
		Name:      name,
		Address:   address,
		City:      city,
		State:     state,
		Phones:    phones,
		CreatedBy: createdBy,
		CreatedAt: dt,
		UpdatedAt: dt,
	}
	var id *string
	id, err = s.repo.Create(newPatient)
	if err != nil {
		return nil, err
	}
	if (id != nil) && (*id != newPatient.ID) {
		return nil, errors.New("New patient creation: returned repository ID not equal to new patient ID")
	}
	return id, err
}

// FindByID searches a patient by its ID
func (s *patientService) FindByID(id string) (*patient.Patient, error) {
	if id == "" {
		return nil, patient.ErrInvalidInputSyntax
	}
	_, err := uuid.FromString(id)
	if err != nil {
		return nil, patient.ErrInvalidInputSyntax
	}

	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateName, updates patient name
func (s *patientService) UpdateName(id string, name, updatedBy string) error {
	_, err := uuid.FromString(id)
	if err != nil {
		return patient.ErrInvalidInputSyntax
	}
	_, err = uuid.FromString(updatedBy)
	if err != nil {
		return patient.ErrInvalidInputSyntax
	}
	if name == "" {
		return patient.ErrInvalidInputSyntax
	}
	return s.repo.UpdateName(id, name, updatedBy)
}

// UpdatePhone updates patient email
func (s *patientService) UpdatePhone(id string, phones []string,
	updatedBy string) error {
	_, err := uuid.FromString(id)
	if err != nil {
		return patient.ErrInvalidInputSyntax
	}
	_, err = uuid.FromString(updatedBy)
	if err != nil {
		return patient.ErrInvalidInputSyntax
	}
	if len(phones) > 10 ||
		len(phones) == 0 {
		return patient.ErrInvalidInputSyntax
	}
	if strings.TrimSpace(strings.Join(phones, "")) == "" {
		return patient.ErrInvalidInputSyntax
	}
	return s.repo.UpdatePhone(id, phones, updatedBy)
}

// GetAll returns a paginated list of all patients ordered byname
func (s *patientService) GetAll(previous string, next string, pgSize int) (*patient.Page, error) {
	var page patient.Page
	var err error
	var list *[]patient.Patient
	if pgSize <= 0 {
		return nil, patient.ErrInvalidInputSyntax
	}
	switch {
	case (previous != "" && next != ""):
		// if both (previous & next) are present, return error
		return nil, user.ErrInvalidInputSyntax
	case (previous == "" && next == ""):
		// if both are empty, get the first "pgSize" elements of the list
		list, page.HasNextPage, err = s.repo.
			GetAll("", true, pgSize)
		if err != nil {
			return nil, err
		}
		if list != nil && len(*list) > 0 {
			for _, item := range *list {
				page.Patients = append(page.Patients, item)
			}
		}
		if len(page.Patients) > 0 {
			page.StartCursor = base64.StdEncoding.
				EncodeToString([]byte(page.Patients[0].ID))
			page.EndCursor = base64.StdEncoding.
				EncodeToString([]byte(page.
					Patients[len(page.Patients)-1].ID))
		} else {
			page.StartCursor = ""
			page.EndCursor = ""
		}
		page.HasPreviousPage = false
	case (previous != ""):
		// if previous is present, get a previous list
		c, err := base64.StdEncoding.DecodeString(previous)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		list, page.HasPreviousPage, err = s.repo.GetAll(cursor, false, pgSize)
		if err != nil {
			//fmt.Println(err)
			return nil, err
		}
		if list != nil && len(*list) > 0 {
			for _, item := range *list {
				page.Patients = append(page.Patients, item)
			}
		}
		if len(page.Patients) > 0 {
			//fmt.Printf("StartCursor: %v\n", page.Patients[0].ID)
			page.StartCursor = base64.StdEncoding.
				EncodeToString([]byte(page.Patients[0].ID))
			page.EndCursor = base64.StdEncoding.
				EncodeToString([]byte(page.
					Patients[len(page.Patients)-1].ID))
			page.HasNextPage = true
		} else {
			page.StartCursor = ""
			page.EndCursor = ""
			page.HasNextPage = false
		}
	case (next != ""):
		c, err := base64.StdEncoding.DecodeString(next)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		list, page.HasNextPage, err = s.repo.
			GetAll(cursor, true, pgSize)
		if err != nil {
			return nil, err
		}
		if list != nil && len(*list) > 0 {
			for _, item := range *list {
				page.Patients = append(page.Patients, item)
			}
		}
		if len(page.Patients) > 0 {
			page.StartCursor = base64.StdEncoding.
				EncodeToString([]byte(page.Patients[0].ID))
			page.EndCursor = base64.StdEncoding.
				EncodeToString([]byte(page.
					Patients[len(page.Patients)-1].ID))
			page.HasPreviousPage = true
		} else {
			page.StartCursor = ""
			page.EndCursor = ""
			page.HasPreviousPage = false
		}
	}
	return &page, nil
}

// FindByName returns a list of patients whose names looks like 'name'
func (s *patientService) FindByName(name string) (*[]patient.Patient, error) {
	if len(name) == 0 {
		return nil, patient.ErrInvalidInputSyntax
	}
	res, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	return res, nil
}
