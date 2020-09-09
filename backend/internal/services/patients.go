package services

import (
	"encoding/base64"
	"errors"
	"time"

	"carlosapgomes.com/sked/internal/patient"
	uuid "github.com/satori/go.uuid"
)

// UserService provides implementation of user domain interface
type patientService struct {
	repo patient.Repository
}

// NewUserService returns an user Service instance
func NewPatientService(repo patient.Repository) patient.Service {
	return &patientService{
		repo,
	}
}

// Create creates a new patient and returns its uuid
func (s *patientService) Create(name, address, city, state string, phone []string, createdBy string) (*string, error) {
	uid := uuid.NewV4()
	dt := time.Now().UTC()
	newPatient := patient.Patient{
		ID:        uid.String(),
		Name:      name,
		Address:   address,
		City:      city,
		State:     state,
		Phones:    phone,
		CreatedBy: createdBy,
		CreatedAt: dt,
		UpdatedAt: dt,
	}
	var id *string
	id, err := s.repo.Create(newPatient)
	if err != nil {
		return nil, err
	}
	if (id != nil) && (*id != newPatient.ID) {
		return nil, errors.New("New user creation: returned repository ID not equal to new user ID")
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

// UpdateName, updates user name
func (s *patientService) UpdateName(id string, name string) error {
	if name == "" {
		return patient.ErrInvalidInputSyntax
	}
	return s.repo.UpdateName(id, name)
}

func (s *patientService) GetAll(before string, after string, pgSize int) (*patient.Cursor, error) {
	var patientsResp patient.Cursor
	var err error
	var pList *[]patient.Patient
	switch {
	case (before != "" && after != ""):
		// if both (before & after) are present, returns error
		return nil, patient.ErrInvalidInputSyntax
	case (before == "" && after == ""):
		// if they are empty/or absent
		// get default list and page size
		pList, patientsResp.HasBefore, err = s.repo.
			GetAll("", false, pgSize)
		if err != nil {
			return nil, err
		}
		if pList != nil {
			for _, u := range *pList {
				patientsResp.Patients = append(patientsResp.Patients, u)
			}
		}
		if len(patientsResp.Patients) > 0 {
			patientsResp.Before = base64.StdEncoding.
				EncodeToString([]byte(patientsResp.Patients[len(patientsResp.Patients)-1].ID))
		} else {
			patientsResp.Before = ""
		}
		patientsResp.After = ""
		patientsResp.HasAfter = false
		// and return values
	case (before != ""):
		// if before is present,
		// get a before list
		c, err := base64.StdEncoding.DecodeString(before)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		pList, patientsResp.HasBefore, err = s.repo.GetAll(cursor, false, pgSize)
		if err != nil {
			return nil, err
		}
		if pList != nil {
			for _, u := range *pList {
				patientsResp.Patients = append(patientsResp.Patients, u)
			}
		}
		if len(patientsResp.Patients) > 0 {
			befCursor := base64.StdEncoding.EncodeToString([]byte(patientsResp.Patients[len(patientsResp.Patients)-1].ID))
			patientsResp.Before = befCursor
		} else {
			patientsResp.Before = ""
		}
		// test for 'after data' from the requested cursor
		// fill the response fields
		_, patientsResp.HasAfter, err = s.repo.GetAll(cursor, true, pgSize)
		if patientsResp.HasAfter {
			patientsResp.After = base64.StdEncoding.EncodeToString([]byte(before))
		} else {
			patientsResp.After = ""
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
		pList, patientsResp.HasAfter, err = s.repo.
			GetAll(cursor, true, pgSize)
		// and return it
		if pList != nil {
			for _, p := range *pList {
				patientsResp.Patients = append(patientsResp.Patients, p)
			}
		}
		if len(patientsResp.Patients) > 0 {
			patientsResp.After = base64.StdEncoding.EncodeToString([]byte(patientsResp.Patients[0].Name))
		}
		// test for 'before data' from the requested cursor
		// fill the response fields
		_, patientsResp.HasBefore, err = s.repo.
			GetAll(cursor, false, pgSize)
		if patientsResp.HasBefore {
			patientsResp.Before = base64.StdEncoding.EncodeToString([]byte(after))
		}
	}
	return &patientsResp, nil
}

// FindByName returns a list of users whose names looks like 'name'
func (s *patientService) FindByName(name string) (*[]patient.Patient, error) {
	res, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	return res, nil
}
