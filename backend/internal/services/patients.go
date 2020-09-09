package services

import (
	"encoding/base64"
	"errors"
	"time"

	"carlosapgomes.com/sked/internal/patient"
	"carlosapgomes.com/sked/internal/user"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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
func (s *patientService) FindByID(uID string) (*patient.Patient, error) {
	if uID == "" {
		return nil, user.ErrInvalidInputSyntax
	}
	_, err := uuid.FromString(uID)
	if err != nil {
		return nil, user.ErrInvalidInputSyntax
	}

	u, err := s.repo.FindByID(uID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail searches a user by its email address
func (s *patientService) FindByEmail(email string) (*patient.Patient, error) {
	if email == "" {
		return nil, user.ErrInvalidInputSyntax
	}

	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Authenticate autheticates a user and return its ID
func (s *patientService) Authenticate(email, password string) (*string, error) {
	if (email == "") || (password == "") {
		return nil, user.ErrInvalidInputSyntax
	}
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(u.HashedPw, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, user.ErrInvalidCredentials
		}
		return nil, err
	}
	return &u.ID, nil
}

// UpdateRoles set a new user roles
func (s *patientService) UpdateRoles(id string, roles []string) error {
	return s.repo.UpdateRoles(id, roles)
}

// UpdatePw updates a user's password
func (s *patientService) UpdatePw(id, pw string) error {
	if pw == "" {
		return user.ErrInvalidInputSyntax
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	if err != nil {
		return err
	}
	return s.repo.UpdatePw(id, hashedPw)
}

// UpdateStatus updates a user's status
func (s *patientService) UpdateStatus(id string, active bool) error {
	return s.repo.UpdateStatus(id, active)
}

// UpdateEmailValidated updates a user's EmailValidated
func (s *patientService) UpdateEmailValidated(id string, emailWasValidated bool) error {
	return s.repo.UpdateEmailValidated(id, emailWasValidated)
}

// UpdateName, updates user name
func (s *patientService) UpdateName(id string, name string) error {
	if name == "" {
		return user.ErrInvalidInputSyntax
	}
	return s.repo.UpdateName(id, name)
}

// UpdateEmail, updates user email, only admin should be able to do this.
// email account is the source of truth for system access
func (s *patientService) UpdateEmail(id string, email string) error {
	// validate email string
	if email == "" {
		return user.ErrInvalidInputSyntax
	}
	return s.repo.UpdateEmail(id, email)
}

// UpdatePhone updates user email
func (s *patientService) UpdatePhone(id string, phone string) error {
	if phone == "" {
		return user.ErrInvalidInputSyntax
	}
	return s.repo.UpdatePhone(id, phone)
}

// GetAll returns a paginated list of all users ordered by email
func (s *patientService) GetAll(before string, after string, pgSize int) (*patient.Cursor, error) {
	var patientsResp patient.Cursor
	var err error
	var uList *[]patient.Patient
	switch {
	case (before != "" && after != ""):
		// if both (before & after) are present, returns error
		return nil, user.ErrInvalidInputSyntax
	case (before == "" && after == ""):
		// if they are empty/or absent
		// get default list and page size
		uList, patientsResp.HasBefore, err = s.repo.
			GetAll("", false, pgSize)
		if err != nil {
			return nil, err
		}
		if uList != nil {
			for _, u := range *uList {
				patientsResp.Patients = append(patientsResp.Patients, u)
			}
		}
		if len(patientsResp.Patients) > 0 {
			patientsResp.Before = base64.StdEncoding.
				EncodeToString([]byte(patientsResp.Patients[len(patientsResp.Patients)-1].Email))
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
		uList, patientsResp.HasBefore, err = s.repo.GetAll(cursor, false, pgSize)
		if err != nil {
			return nil, err
		}
		if uList != nil {
			for _, u := range *uList {
				patientsResp.Patients = append(patientsResp.Patients, u)
			}
		}
		if len(patientsResp.Patients) > 0 {
			befCursor := base64.StdEncoding.EncodeToString([]byte(patientsResp.Patients[len(patientsResp.Patients)-1].Email))
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
		uList, patientsResp.HasAfter, err = s.repo.
			GetAll(cursor, true, pgSize)
		// and return it
		if uList != nil {
			for _, u := range *uList {
				patientsResp.Patients = append(patientsResp.Patients, u)
			}
		}
		if len(patientsResp.Patients) > 0 {
			patientsResp.After = base64.StdEncoding.EncodeToString([]byte(patientsResp.Patients[0].Email))
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
