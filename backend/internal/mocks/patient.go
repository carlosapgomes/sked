package mocks

import (
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/patient"
	"github.com/pkg/errors"
)

// mocks user services and repository

// PatientMockRepo is a mocked user repository
type PatientMockRepo struct {
	pDb []patient.Patient
}

//NewPatientRepo returns a mocked repository
func NewPatientRepo() *PatientMockRepo {
	var db []patient.Patient
	validPatient := &patient.Patient{
		ID:        "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		Name:      "Valid Patient",
		Phones:    []string{"6544332135"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	validPatient2 := &patient.Patient{
		ID:        "68b1d5e2-39dd-4713-8631-a08100383a0f",
		Name:      "Bob",
		Phones:    []string{"6544334535"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	validPatient3 := &patient.Patient{
		ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
		Name:      "Alice Jones",
		Phones:    []string{"6544332135"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	validPatient4 := &patient.Patient{
		ID:        "ecadbb28-14e6-4560-8574-809c6c54b9cb",
		Name:      "Barack Obama",
		Phones:    []string{"6544332135"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	validPatient5 := &patient.Patient{
		ID:        "ca16fc9d-df7b-4594-97e3-264432145b01",
		Name:      "SpongeBob Squarepants",
		Phones:    []string{"65949340"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	validPatient6 := &patient.Patient{
		ID:        "27f9802b-acb3-4852-bf97-c4ed4c3b3658",
		Name:      "Tim Berners-Lee",
		Phones:    []string{"0323949324"},
		CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	db = append(db, *validPatient)
	db = append(db, *validPatient2)
	db = append(db, *validPatient3)
	db = append(db, *validPatient4)
	db = append(db, *validPatient5)
	db = append(db, *validPatient6)
	return &PatientMockRepo{
		db,
	}
}

// Create mocks user creation
func (r *PatientMockRepo) Create(user patient.Patient) (*string, error) {
	badID := "12342342"
	switch user.Name {
	case "Valid patient":
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
	for i, u := range r.pDb {
		if u.ID == id {
			r.pDb[i].Name = name
			return nil
		}
	}
	return patient.ErrNoRecord
}

// UpdatePhone mocks updating user's Phone
func (r *PatientMockRepo) UpdatePhone(id string, phones []string) error {
	for i, u := range r.pDb {
		if u.ID == id {
			r.pDb[i].Phones = phones
			return nil
		}
	}
	return patient.ErrNoRecord
}

// FindByID mocks finding a user by its ID
func (r *PatientMockRepo) FindByID(id string) (*patient.Patient, error) {
	for _, u := range r.pDb {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, patient.ErrNoRecord
}

// GetAll returns a list of users ordered by email
func (r *PatientMockRepo) GetAll(cursor string, after bool, pgSize int) (*[]patient.Patient, bool, error) {
	if cursor == "" {
		return &r.pDb, false, nil
	}
	pos := r.findPos(r.pDb, cursor)
	if pos == -1 {
		return nil, false, patient.ErrNoRecord
	}

	var res []patient.Patient
	var hasMore bool
	hasMore = false
	if after {
		start := pos + 1
		for i := start; i < (start + pgSize); i++ {
			res = append(res, r.pDb[i])
		}
		if (len(r.pDb) - pos) > pgSize {
			hasMore = true
		}
	} else {
		start := pos - pgSize
		if start < 0 {
			start = 0
		}
		for i := start; i <= (pos - 1); i++ {
			res = append(res, r.pDb[i])
		}
		if pos > pgSize {
			hasMore = true
		}
	}
	return &res, hasMore, nil
}

func (r PatientMockRepo) findPos(patients []patient.Patient, id string) int {
	for i, el := range patients {
		if el.ID == id {
			return i
		}
	}
	return -1
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
func (s PatientMockSvc) UpdatePhone(id string, phones []string) error {
	return nil
}

// GetAll return a lista of users ordered by email
func (s PatientMockSvc) GetAll(before string, after string, pgSize int) (*patient.Cursor, error) {
	return nil, nil
}

// FindByName return a list of users whose names looks like 'name'
func (s PatientMockSvc) FindByName(name string) (*[]patient.Patient, error) {
	repo := NewPatientRepo()
	var users []patient.Patient
	for _, u := range repo.pDb {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(name)) {
			users = append(users, u)
		}
	}
	if len(users) == 0 {
		return nil, patient.ErrNoRecord
	}
	return &users, nil
}
