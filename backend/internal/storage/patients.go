package storage

import (
	"database/sql"

	"carlosapgomes.com/sked/internal/patient"
)

// patientRepository type
type patientRepository struct {
	DB *sql.DB
}

func NewPgPatientRepository(db *sql.DB) patient.Repository {
	return &patientRepository{
		db,
	}
}

// Create - creates a new patient record
func (r patientRepository) Create(p patient.Patient) (*string, error) {
	return nil, nil
}

// UpdateName - updates a patient's name
func (r patientRepository) UpdateName(id, name string) error {
	return nil
}

// UpdatePhone - update a patient's phones list
func (r patientRepository) UpdatePhone(id string, phones []string) error {
	return nil
}

// FindByID - finds a patient by its ID
func (r patientRepository) FindByID(id string) (*patient.Patient, error) {
	return nil, nil
}

// FindByName - find a patient by its name
func (r patientRepository) FindByName(name string) (*[]patient.Patient, error) {
	return nil, nil
}

// GetAll - returns a paginated list of patients
func (r patientRepository) GetAll(cursor string, after bool,
	pgSize int) (*[]patient.Patient, bool, error) {
	return nil, false, nil
}
