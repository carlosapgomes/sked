package storage

import (
	"database/sql"
	"time"

	"carlosapgomes.com/sked/internal/surgery"
)

// surgeryRepository type
type surgeryRepository struct {
	DB *sql.DB
}

// NewPgSurgeryRepository returns an instance of an surgeryRepository
func NewPgSurgeryRepository(db *sql.DB) surgery.Repository {
	return &surgeryRepository{
		db,
	}
}

// Create - creates a new surgery
func (r surgeryRepository) Create(surg surgery.Surgery) (*string, error) {
	return nil, nil
}

// Update - updates an surgery
func (r surgeryRepository) Update(surg surgery.Surgery) (*string, error) {
	return nil, nil
}

// FindByID - finds an surgery by its ID
func (r surgeryRepository) FindByID(id string) (*surgery.Surgery, error) {
	return nil, nil
}

// FindByPatientID - finds a surgery by its patientID
func (r surgeryRepository) FindByPatientID(patientID string) ([]*surgery.Surgery, error) {
	return nil, nil
}

// FindFindByDoctorID - finds a surgery by its doctorID
func (r surgeryRepository) FindByDoctorID(doctorID string) ([]*surgery.Surgery, error) {
	return nil, nil
}

// FindByDate - finds surgeries in a date
func (r surgeryRepository) FindByDate(date time.Time) ([]*surgery.Surgery, error) {
	return nil, nil
}

// GetAll - returns a paginated list of surgeries
func (r surgeryRepository) GetAll(cursor string, after bool, pgSize int) (*[]surgery.Surgery, bool, error) {
	return nil, false, nil
}
