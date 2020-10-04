package storage

import (
	"database/sql"
	"fmt"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
	"carlosapgomes.com/sked/internal/user"
	"github.com/lib/pq"
)

// appointmentRepository type
type appointmentRepository struct {
	DB *sql.DB
}

// NewPgAppointmentRepository returns an instance of an appointmentRepository
func NewPgAppointmentRepository(db *sql.DB) appointment.Repository {
	return &appointmentRepository{
		db,
	}
}

// Create - creates a new appointment
func (r appointmentRepository) Create(a appointment.Appointment) (*string, error) {
	stmt := `INSERT INTO appointments (id, date_time, patient_name,
		patient_id, doctor_name, doctor_id, notes, canceled, completed,
		created_by, created_at, updated_by, updated_at) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) Returning id;`
	var id string
	err := r.DB.QueryRow(stmt, a.ID, a.DateTime, a.PatientName, a.PatientID,
		a.DoctorName, a.DoctorID, a.Notes, a.Canceled, a.Completed,
		a.CreatedBy, a.CreatedAt, a.UpdatedBy, a.UpdatedAt).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return nil, fmt.Errorf("%w\n %s %s", appointment.ErrDuplicateField,
				pqErr.Message, pqErr.Column)
		case "22P02":
			return nil, fmt.Errorf("%w\n %s %s", appointment.ErrInvalidInputSyntax,
				pqErr.Message, pqErr.Column)
		}
		return nil, fmt.Errorf("%w\n %s %s", appointment.ErrDb, pqErr.Message,
			pqErr.Column)
	}
	return &id, err
}

// Update - updates an appointment
func (r appointmentRepository) Update(a appointment.Appointment) (*string, error) {
	stmt := `UPDATE appointments SET date_time = $1, notes = $2, canceled = $3,
		completed = $4, updated_by = $5, updated_at = %6 WHERE id = 7`
	_, err := r.DB.Exec(stmt, a.DateTime, a.Notes, a.Canceled, a.Completed,
		a.UpdatedBy, a.UpdatedAt, a.ID)
	if err != nil {
		pqErr := err.(*pq.Error)
		return nil, pqErr
	}
	return &a.ID, err
}

// FindByID - finds an appointment by its ID
func (r appointmentRepository) FindByID(id string) (*appointment.Appointment, error) {
	var a appointment.Appointment
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name, doctor_id,
			notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID, &a.DoctorName,
		&a.DoctorID, &a.Notes, &a.Canceled, &a.Completed, &a.CreatedBy,
		&a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, user.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	a.CreatedAt = a.CreatedAt.In(loc)
	a.UpdatedAt = a.UpdatedAt.In(loc)
	return &a, err
}

// FindByPatientID - finds an appointment by its patientID
func (r appointmentRepository) FindByPatientID(patientID string) ([]appointment.Appointment, error) {
	var apptmts []appointment.Appointment
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name, doctor_id,
			notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at WHERE patient_id = $1`
	rows, err := r.DB.Query(stmt, patientID)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var a appointment.Appointment
		err = rows.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID, &a.DoctorName,
			&a.DoctorID, &a.Notes, &a.Canceled, &a.Completed, &a.CreatedBy,
			&a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, user.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		a.CreatedAt = a.CreatedAt.In(loc)
		a.UpdatedAt = a.UpdatedAt.In(loc)
		apptmts = append(apptmts, a)
	}
	// every date/time was saved as UTC, so use them as UTC
	return &apptmts, err
}

// FindFindByDoctorID - finds an appointment by its doctorID
func (r appointmentRepository) FindByDoctorID(doctorID string) ([]*appointment.Appointment, error) {
	var a appointment.Appointment
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name, doctor_id,
			notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at WHERE doctor_id = $1`
	row := r.DB.QueryRow(stmt, doctorID)
	err := row.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID, &a.DoctorName,
		&a.DoctorID, &a.Notes, &a.Canceled, &a.Completed, &a.CreatedBy,
		&a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, user.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	a.CreatedAt = a.CreatedAt.In(loc)
	a.UpdatedAt = a.UpdatedAt.In(loc)
	return &a, err
}

// FindByDate - finds appointments in a date
func (r appointmentRepository) FindByDate(date time.Time) ([]*appointment.Appointment, error) {
	return nil, nil
}

// GetAll - returns a paginated list of appointments
func (r appointmentRepository) GetAll(cursor string, after bool, pgSize int) (*[]appointment.Appointment, bool, error) {
	return nil, false, nil
}
