package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
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
	err := r.DB.QueryRow(stmt,
		a.ID,
		strings.TrimRight(a.DateTime.String(), " UTC"),
		a.PatientName,
		a.PatientID,
		a.DoctorName,
		a.DoctorID,
		a.Notes,
		a.Canceled,
		a.Completed,
		a.CreatedBy,
		strings.TrimRight(a.CreatedAt.String(), " UTC"),
		a.UpdatedBy,
		strings.TrimRight(a.UpdatedAt.String(), " UTC")).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok {
		fmt.Printf("db error: %v\n", pqErr)
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
		completed = $4, updated_by = $5, updated_at = $6 WHERE id = $7`
	_, err := r.DB.Exec(stmt, a.DateTime, a.Notes, a.Canceled, a.Completed,
		a.UpdatedBy, a.UpdatedAt, a.ID)
	if err != nil {
		pqErr := err.(*pq.Error)
		return nil, pqErr
	}
	return &a.ID, err
}

// FindByID - finds an appointment by its ID
func (r appointmentRepository) FindByID(id string) (*appointment.Appointment,
	error) {
	var a appointment.Appointment
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at FROM appointments WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID,
		&a.DoctorName, &a.DoctorID, &a.Notes, &a.Canceled, &a.Completed,
		&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, appointment.ErrNoRecord
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
func (r appointmentRepository) FindByPatientID(patientID string) ([]appointment.Appointment,
	error) {
	var apptmts []appointment.Appointment
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name, doctor_id,
			notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at FROM appointments WHERE patient_id = $1`
	rows, err := r.DB.Query(stmt, patientID)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var a appointment.Appointment
		err = rows.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID,
			&a.DoctorName, &a.DoctorID, &a.Notes, &a.Canceled, &a.Completed,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, appointment.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		a.CreatedAt = a.CreatedAt.In(loc)
		a.UpdatedAt = a.UpdatedAt.In(loc)
		apptmts = append(apptmts, a)
	}
	// every date/time was saved as UTC, so use them as UTC
	return apptmts, err
}

// FindFindByDoctorID - finds an appointment by its doctorID
func (r appointmentRepository) FindByDoctorID(doctorID string) ([]appointment.Appointment,
	error) {
	var apptmts []appointment.Appointment
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name, doctor_id,
			notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at FROM appointments WHERE doctor_id = $1`
	rows, err := r.DB.Query(stmt, doctorID)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var a appointment.Appointment
		err = rows.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID,
			&a.DoctorName, &a.DoctorID, &a.Notes, &a.Canceled, &a.Completed,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, appointment.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		a.CreatedAt = a.CreatedAt.In(loc)
		a.UpdatedAt = a.UpdatedAt.In(loc)
		apptmts = append(apptmts, a)
	}
	return apptmts, err
}

// FindByMonthYear - return all appointments in a specific month
func (r appointmentRepository) FindByInterval(s,
	e time.Time) ([]appointment.Appointment, error) {
	start := s.Format("2006-01-02")
	end := e.Format("2006-01-02")
	var apptmts []appointment.Appointment
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at FROM appointments
			WHERE sked_date_to_char(start)  >= $1 AND 
			sked_date_to_char(end)  <= $2`
	rows, err := r.DB.Query(stmt, start, end)
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var a appointment.Appointment
		err = rows.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID,
			&a.DoctorName, &a.DoctorID, &a.Notes, &a.Canceled, &a.Completed,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, appointment.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		a.CreatedAt = a.CreatedAt.In(loc)
		a.UpdatedAt = a.UpdatedAt.In(loc)
		apptmts = append(apptmts, a)
	}
	return apptmts, err
}

// FindByDate - finds appointments in a date
func (r appointmentRepository) FindByDate(date time.Time) ([]appointment.Appointment, error) {
	var apptmts []appointment.Appointment
	searchDate := date.Format("2006-01-02")
	stmt := `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at FROM appointments
			WHERE sked_date_to_char(date_time) = $1`
	rows, err := r.DB.Query(stmt, searchDate)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("UTC")
	defer rows.Close()
	for rows.Next() {
		var a appointment.Appointment
		err = rows.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID,
			&a.DoctorName, &a.DoctorID, &a.Notes, &a.Canceled, &a.Completed,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, appointment.ErrNoRecord
		} else if err != nil {
			return nil, err
		}
		a.CreatedAt = a.CreatedAt.In(loc)
		a.UpdatedAt = a.UpdatedAt.In(loc)
		apptmts = append(apptmts, a)
	}
	return apptmts, err
}

// GetAll - returns a paginated list of appointments
func (r appointmentRepository) GetAll(cursor string, next bool,
	pgSize int) ([]appointment.Appointment, bool, error) {
	if pgSize <= 0 {
		pgSize = 15
	}
	appointmt, err := r.FindByID(cursor)
	if err != nil {
		return nil, false, err
	}

	cursorDate := appointmt.DateTime.Format("2006-01-02")
	fmt.Printf("cursorDate: %v\n", cursorDate)
	var stmt string
	if next {
		// Get next result page
		stmt = `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at from appointments 
			WHERE sked_date_to_char(date_time) > $1
			ORDER BY date_time LIMIT $2`
	} else {
		// Get previous result page
		stmt = `SELECT id, date_time, patient_name, patient_id, doctor_name,
			doctor_id, notes, canceled, completed, created_by, created_at, 
			updated_by, updated_at from appointments 
			WHERE sked_date_to_char(date_time) < $1
			ORDER BY date_time LIMIT $2`
	}
	// request one more item to help set 'hasMore' flag (see bellow)
	rows, err := r.DB.Query(stmt, cursorDate, (pgSize + 1))
	if err != nil {
		return nil, false, err
	}
	var appointmts []appointment.Appointment
	defer rows.Close()
	for rows.Next() {
		var a appointment.Appointment
		err = rows.Scan(&a.ID, &a.DateTime, &a.PatientName, &a.PatientID,
			&a.DoctorName, &a.DoctorID, &a.Notes, &a.Canceled, &a.Completed,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedBy, &a.UpdatedAt)
		if err == sql.ErrNoRows {
			return nil, false, appointment.ErrNoRecord
		} else if err != nil {
			return nil, false, err
		}
		appointmts = append(appointmts, a)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, false, err
	}
	hasMore := false
	if len(appointmts) == (pgSize + 1) {
		// remove the element that was not requested
		if next {
			// remove last element
			appointmts = appointmts[:len(appointmts)-1]
		} else {
			// remove first element
			appointmts = appointmts[1:]
		}
		hasMore = true
	}
	return appointmts, hasMore, nil
}
