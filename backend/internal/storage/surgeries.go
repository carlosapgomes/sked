package storage

import (
	"database/sql"
	"fmt"
	"time"

	"carlosapgomes.com/sked/internal/surgery"
	"github.com/lib/pq"
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
	stmt := `INSERT INTO surgeries (id, date_time, patient_name,
		patient_id, doctor_name, doctor_id, notes, proposed_surgery, canceled, done,
		created_by, created_at, updated_by, updated_at) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) Returning id;`
	var id string
	err := r.DB.QueryRow(stmt, surg.ID, surg.DateTime, surg.PatientName, surg.PatientID,
		surg.DoctorName, surg.DoctorID, surg.Notes, surg.ProposedSurgery,
		surg.Canceled, surg.Done, surg.CreatedBy, surg.CreatedAt, surg.UpdatedBy,
		surg.UpdatedAt).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return nil, fmt.Errorf("%w\n %s %s", surgery.ErrDuplicateField,
				pqErr.Message, pqErr.Column)
		case "22P02":
			return nil, fmt.Errorf("%w\n %s %s", surgery.ErrInvalidInputSyntax,
				pqErr.Message, pqErr.Column)
		}
		return nil, fmt.Errorf("%w\n %s %s", surgery.ErrDb, pqErr.Message,
			pqErr.Column)
	}
	return &id, err
}

// Update - updates an surgery
func (r surgeryRepository) Update(surg surgery.Surgery) (*string, error) {
	stmt := `UPDATE surgeries SET date_time = $1, notes = $2, 
		proposed_surgery = $3, canceled = $4,
		done = $5, updated_by = $6, updated_at = $7 WHERE id = $8`
	_, err := r.DB.Exec(stmt, surg.DateTime, surg.Notes, surg.Canceled, surg.Done,
		surg.UpdatedBy, surg.UpdatedAt, surg.ID)
	if err != nil {
		pqErr := err.(*pq.Error)
		return nil, pqErr
	}
	return &surg.ID, err
}

// FindByID - finds an surgery by its ID
func (r surgeryRepository) FindByID(id string) (*surgery.Surgery, error) {
	var s surgery.Surgery
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, proposed_surgery, canceled, done, 
			created_by, created_at, updated_by, updated_at
			FROM surgeries WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&s.ID, &s.DateTime, &s.PatientName, &s.PatientID,
		&s.DoctorName, &s.DoctorID, &s.Notes, &s.ProposedSurgery, &s.Canceled,
		&s.Done, &s.CreatedBy, &s.CreatedAt, &s.UpdatedBy, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, surgery.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	s.CreatedAt = s.CreatedAt.In(loc)
	s.UpdatedAt = s.UpdatedAt.In(loc)
	return &s, err
}

// FindByPatientID - finds a surgery by its patientID
func (r surgeryRepository) FindByPatientID(patientID string) ([]surgery.Surgery, error) {
	var surgs []surgery.Surgery
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name, doctor_id,
			notes, proposed_surgery, canceled, done, created_by, created_at, 
			updated_by, updated_at FROM surgeries WHERE patient_id = $1`
	rows, err := r.DB.Query(stmt, patientID)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var s surgery.Surgery
		err = rows.Scan(&s.ID, &s.DateTime, &s.PatientName, &s.PatientID,
			&s.DoctorName, &s.DoctorID, &s.Notes, &s.ProposedSurgery,
			&s.Canceled, &s.Done, &s.CreatedBy, &s.CreatedAt,
			&s.UpdatedBy, &s.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, surgery.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		s.CreatedAt = s.CreatedAt.In(loc)
		s.UpdatedAt = s.UpdatedAt.In(loc)
		surgs = append(surgs, s)
	}
	// every date/time was saved as UTC, so use them as UTC
	return surgs, err
}

// FindFindByDoctorID - finds a surgery by its doctorID
func (r surgeryRepository) FindByDoctorID(doctorID string) ([]surgery.Surgery, error) {
	var surgs []surgery.Surgery
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name, doctor_id,
			notes, proposed_surgery, canceled, done, created_by, created_at, 
			updated_by, updated_at FROM surgeries WHERE doctor_id = $1`
	rows, err := r.DB.Query(stmt, doctorID)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var s surgery.Surgery
		err = rows.Scan(&s.ID, &s.DateTime, &s.PatientName, &s.PatientID,
			&s.DoctorName, &s.DoctorID, &s.Notes, &s.ProposedSurgery,
			&s.Canceled, &s.Done, &s.CreatedBy, &s.CreatedAt,
			&s.UpdatedBy, &s.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, surgery.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		s.CreatedAt = s.CreatedAt.In(loc)
		s.UpdatedAt = s.UpdatedAt.In(loc)
		surgs = append(surgs, s)
	}
	return surgs, err
}

// FindByDate - finds surgeries in a date
func (r surgeryRepository) FindByDate(date time.Time) ([]surgery.Surgery, error) {
	var surgs []surgery.Surgery
	searchDate := date.Format("2006-01-02")
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, proposed_surgery, canceled, done,
			created_by, created_at, updated_by, updated_at
			FROM surgeries
			WHERE sked_date_to_char(date_time) = $1`
	rows, err := r.DB.Query(stmt, searchDate)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var s surgery.Surgery
		err = rows.Scan(&s.ID, &s.DateTime, &s.PatientName, &s.PatientID,
			&s.DoctorName, &s.DoctorID, &s.Notes, &s.ProposedSurgery,
			&s.Canceled, &s.Done, &s.CreatedBy, &s.CreatedAt,
			&s.UpdatedBy, &s.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, surgery.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		s.CreatedAt = s.CreatedAt.In(loc)
		s.UpdatedAt = s.UpdatedAt.In(loc)
		surgs = append(surgs, s)
	}
	return surgs, err
}

// GetAll - returns a paginated list of surgeries
func (r surgeryRepository) GetAll(cursor string, next bool, pgSize int) ([]surgery.Surgery, bool, error) {
	if pgSize <= 0 {
		pgSize = 15
	}
	surg, err := r.FindByID(cursor)
	if err != nil {
		return nil, false, err
	}
	cursorDate := surg.DateTime.Format("2006-01-02")
	fmt.Printf("cursorDate: %v\n", cursorDate)
	var stmt string
	if next {
		// Get next result page
		stmt = `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, proposed_surgery, canceled, done, 
			created_by, created_at, updated_by, updated_at from surgeries 
			WHERE sked_date_to_char(date_time) > $1
			ORDER BY date_time LIMIT $2`
	} else {
		// Get previous result page
		stmt = `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, proposed_surgery, canceled, done, 
			created_by, created_at, updated_by, updated_at from surgeries 
			WHERE sked_date_to_char(date_time) < $1
			ORDER BY date_time LIMIT $2`
	}
	// request one more item to help set 'hasMore' flag (see bellow)
	rows, err := r.DB.Query(stmt, cursorDate, (pgSize + 1))
	if err != nil {
		return nil, false, err
	}
	var surgs []surgery.Surgery
	defer rows.Close()
	for rows.Next() {
		var s surgery.Surgery
		err = rows.Scan(&s.ID, &s.DateTime, &s.PatientName, &s.PatientID,
			&s.DoctorName, &s.DoctorID, &s.Notes, &s.ProposedSurgery,
			&s.Canceled, &s.Done, &s.CreatedBy, &s.CreatedAt,
			&s.UpdatedBy, &s.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, false, surgery.ErrNoRecord
		} else if err != nil {
			return nil, false, err
		}
		surgs = append(surgs, s)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, false, err
	}
	hasMore := false
	if len(surgs) == (pgSize + 1) {
		// remove last slice item, because it was not requested
		surgs = surgs[:len(surgs)-1]
		hasMore = true
	}
	return surgs, hasMore, nil
}
