package mocks

import (
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/patient"
	"carlosapgomes.com/sked/internal/user"
	"github.com/pkg/errors"
)

// mocks user services and repository

// PatientMockRepo is a mocked user repository
type PatientMockRepo struct {
	uDb []patient.Patient
}

//NewPatientRepo returns a mocked repository
func NewPatientRepo() *PatientMockRepo {
	var db []patient.Patient
	validUser := &patient.Patient{
		ID:        "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		Name:      "Valid Patient",
		Phones:    []string{"6544332135"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	validUser2 := &patient.Patient{
		ID:        "68b1d5e2-39dd-4713-8631-a08100383a0f",
		Name:      "Bob",
		Phones:    []string{"6544334535"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	validUser3 := &patient.Patient{
		ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
		Name:      "Alice Jones",
		Phones:    []string{"6544332135"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	validUser4 := &patient.Patient{
		ID:        "ecadbb28-14e6-4560-8574-809c6c54b9cb",
		Name:      "Barack Obama",
		Phones:    []string{"6544332135"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	validUser5 := &patient.Patient{
		ID:        "ca16fc9d-df7b-4594-97e3-264432145b01",
		Name:      "SpongeBob Squarepants",
		Phones:    []string{"65949340"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	validUser6 := &patient.Patient{
		ID:        "27f9802b-acb3-4852-bf97-c4ed4c3b3658",
		Name:      "Tim Berners-Lee",
		Phones:    []string{"0323949324"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	db = append(db, *validUser)
	db = append(db, *validUser2)
	db = append(db, *validUser3)
	db = append(db, *validUser4)
	db = append(db, *validUser5)
	db = append(db, *validUser6)
	return &PatientMockRepo{
		db,
	}
}

// Create mocks user creation
func (r *PatientMockRepo) Create(user patient.Patient) (*string, error) {
	badID := "12342342"
	switch user.Name {
	case "New User":
		return &user.ID, nil
	case "Bad uuid":
		return &badID, nil
	case "DB error":
		return nil, errors.New("DB error")
	default:
		return nil, nil
	}
}

// UpdateName mocks updating user's Name
func (r *PatientMockRepo) UpdateName(id string, name string) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].Name = name
			return nil
		}
	}
	return patient.ErrNoRecord
}

// UpdatePhone mocks updating user's Phone
func (r *PatientMockRepo) UpdatePhone(id string, phone string) error {
	var newPhone []string
	newPhone = append(newPhone, phone)
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].Phones = newPhone
			return nil
		}
	}
	return patient.ErrNoRecord
}

// FindByID mocks finding a user by its ID
func (r *PatientMockRepo) FindByID(id string) (*patient.Patient, error) {
	for _, u := range r.uDb {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, patient.ErrNoRecord
}

// GetAll returns a list of users ordered by email
func (r *PatientMockRepo) GetAll(cursor string, after bool, pgSize int) (*[]patient.Patient, bool, error) {
	return nil, false, nil
}

// FindByName returns a list of users whose names looks like 'name'
func (r *PatientMockRepo) FindByName(name string) (*[]patient.Patient, error) {
	return nil, nil
}

// PatientMockSvc mocks user services
type PatientMockSvc struct {
}

// NewPatientSvc returns a mocked user service
func NewPatientSvc() *PatientMockSvc {
	return &PatientMockSvc{}
}

// Create mocks new user creation service
func (s PatientMockSvc) Create(name, address, city, state string, phone []string, createdBy string) (*string, error) {
	switch name {
	case "Bob":
		return nil, errors.Wrap(patient.ErrDuplicateField, "there is already a patient with that name in the database")
	default:
		uid := "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd"
		return &uid, nil
	}
}

// FindByID mocks finding a user by its ID
func (s PatientMockSvc) FindByID(id string) (*patient.Patient, error) {
	if id == "68b1d5e2-39dd-4713-8631-a08100383a0f" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &patient.Patient{
			ID:        "68b1d5e2-39dd-4713-8631-a08100383a0f",
			Name:      "Bob",
			Phones:    []string{"6544334535"},
			CreatedAt: dt,
			UpdatedAt: dt,
		}, nil
	}
	if id == "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &patient.Patient{
			ID:        "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			Name:      "Valid User",
			Phones:    []string{"6544332135"},
			CreatedAt: dt,
			UpdatedAt: dt,
		}, nil
	}
	if id == "dcce1beb-aee6-4a4d-b724-94d470817323" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &patient.Patient{
			ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
			Name:      "Alice Jones",
			Phones:    []string{"6544332135"},
			CreatedAt: dt,
			UpdatedAt: dt,
		}, nil
	}
	if id == "ecadbb28-14e6-4560-8574-809c6c54b9cb" {
		return &patient.Patient{
			ID:        "ecadbb28-14e6-4560-8574-809c6c54b9cb",
			Name:      "Barack Obama",
			Phones:    []string{"6544332135"},
			CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil
	}
	return nil, errors.New("not found")
}

// UpdateName updates user name
func (s PatientMockSvc) UpdateName(id string, name string) error {
	return nil
}

// UpdatePhone updates user email
func (s PatientMockSvc) UpdatePhone(id string, phone string) error {
	return nil
}

// GetAll return a lista of users ordered by email
func (s PatientMockSvc) GetAll(before string, after string, pgSize int) (*user.Cursor, error) {
	return nil, nil
}

// FindByName return a list of users whose names looks like 'name'
func (s PatientMockSvc) FindByName(name string) (*[]patient.Patient, error) {
	repo := NewPatientRepo()
	var users []patient.Patient
	for _, u := range repo.uDb {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(name)) {
			users = append(users, u)
		}
	}
	if len(users) == 0 {
		return nil, patient.ErrNoRecord
	}
	return &users, nil
}
