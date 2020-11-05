package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/patient"
	"github.com/lib/pq"
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
	stmt := `INSERT INTO patients (id, name, address, city, state, phones, 
             created_by, created_at, updated_by, updated_at) VALUES($1,
			 $2, $3, $4, $5, $6, $7, $8, $9, $10) Returning id;`
	var id string
	err := r.DB.QueryRow(
		stmt,
		p.ID,
		p.Name,
		p.Address,
		p.City,
		p.State,
		pq.Array(p.Phones),
		p.CreatedBy,
		strings.TrimRight(p.CreatedAt.String(), " UTC"),
		p.UpdatedBy,
		strings.TrimRight(p.UpdatedAt.String(), " UTC")).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return nil, fmt.Errorf("%w\n %s %s", patient.ErrDuplicateField,
				pqErr.Message, pqErr.Column)
		case "22P02":
			return nil, fmt.Errorf("%w\n %s %s", patient.ErrInvalidInputSyntax,
				pqErr.Message, pqErr.Column)
		}
		return nil, fmt.Errorf("%w\n %s %s", patient.ErrDb, pqErr.Message,
			pqErr.Column)
	}
	return &id, err
}

// UpdateName - updates a patient's name
func (r patientRepository) UpdateName(id, name, updatedBy string) error {
	stmt := `UPDATE patients SET name = $1, updated_by = $2, updated_at = $3
			WHERE id = $4`
	_, err := r.DB.Exec(stmt, name, updatedBy, time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// UpdatePhone - update a patient's phones list
func (r patientRepository) UpdatePhone(id string,
	phones []string, updatedBy string) error {
	stmt := `UPDATE patients SET phones = $1, updated_by = $2, updated_at = $3
			WHERE id = $4`
	_, err := r.DB.Exec(stmt, pq.Array(phones), updatedBy,
		time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// FindByID - finds a patient by its ID
func (r patientRepository) FindByID(id string) (*patient.Patient, error) {
	var p patient.Patient
	stmt := `SELECT name, address, city, state, phones, created_by, 
			created_at, updated_by, updated_at FROM patients WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&p.Name, &p.Address, &p.City, &p.State,
		pq.Array(&p.Phones), &p.CreatedBy, &p.CreatedAt,
		&p.UpdatedBy, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, patient.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	p.ID = id
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	p.CreatedAt = p.CreatedAt.In(loc)
	p.UpdatedAt = p.UpdatedAt.In(loc)
	return &p, err
}

// FindByName - find a patient by its name
func (r patientRepository) FindByName(name string) (*[]patient.Patient, error) {
	// max result size
	maxLstSize := 50
	var patients []patient.Patient
	stmt := `SELECT id, name, address, city, state, phones, created_by, 
			created_at, updated_by, updated_at FROM patients 
			WHERE name ILIKE $1 ORDER BY name LIMIT $2`
	var pattrn strings.Builder
	pattrn.WriteString("%")
	pattrn.WriteString(name)
	pattrn.WriteString("%")
	rows, err := r.DB.Query(stmt, pattrn.String(), maxLstSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p patient.Patient
		err = rows.Scan(&p.ID, &p.Name, &p.Address, &p.City, &p.State,
			pq.Array(&p.Phones), &p.CreatedBy, &p.CreatedAt,
			&p.UpdatedBy, &p.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, patient.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &patients, err
}

// GetAll - returns a paginated list of patients
func (r patientRepository) GetAll(cursor string, next bool,
	pgSize int) (*[]patient.Patient, bool, error) {
	if pgSize <= 0 {
		pgSize = 15
	}
	var stmt string
	var rows *sql.Rows
	var err error
	if cursor == "" {
		stmt = `SELECT id, name, address, city, state, phones, created_by, 
			created_at, updated_by, updated_at FROM patients 
			ORDER BY name LIMIT $1`
		rows, err = r.DB.Query(stmt, (pgSize + 1))
		if err != nil {
			return nil, false, err
		}
	} else {
		p, err := r.FindByID(cursor)
		if err != nil {
			return nil, false, err
		}
		if next {
			// Get next result page
			stmt = `SELECT id, name, address, city, state, phones, created_by, 
			created_at, updated_by, updated_at FROM patients 
			WHERE name > $1 ORDER BY name LIMIT $2`
		} else {
			// Get previous result page
			stmt = `SELECT id, name, address, city, state, phones, created_by, 
			created_at, updated_by, updated_at FROM patients 
			WHERE  name < $1 ORDER BY name LIMIT $2`
		}
		rows, err = r.DB.Query(stmt, p.Name, (pgSize + 1))
		if err != nil {
			return nil, false, err
		}
	}

	var patients []patient.Patient
	defer rows.Close()
	for rows.Next() {
		var p patient.Patient
		err = rows.Scan(&p.ID, &p.Name, &p.Address, &p.City, &p.State,
			pq.Array(&p.Phones), &p.CreatedBy, &p.CreatedAt,
			&p.UpdatedBy, &p.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, false, patient.ErrNoRecord
		} else if err != nil {
			return nil, false, err
		}
		patients = append(patients, p)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, false, err
	}
	hasMore := false
	if len(patients) == (pgSize + 1) {
		// remove the element that was not requested
		if next {
			// remove last element
			patients = patients[:len(patients)-1]
		} else {
			// remove first element
			patients = patients[1:]
		}
		hasMore = true
	}
	return &patients, hasMore, nil
}
