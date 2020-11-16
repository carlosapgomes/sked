package mocks

import (
	"time"

	"carlosapgomes.com/sked/internal/appointment"
)

// mocks appointment services and repository

// AppointmentMockRepo is a mocked appointment repository
type AppointmentMockRepo struct {
	aDb []appointment.Appointment
}

//NewAppointmentRepo returns a mocked repository
func NewAppointmentRepo() *AppointmentMockRepo {
	var db []appointment.Appointment
	return &AppointmentMockRepo{
		db,
	}
}

// Create
func (r AppointmentMockRepo) Create(appointment appointment.Appointment) (*string, error) {
	var id string
	id = appointment.ID
	return &id, nil
}

// Update
func (r AppointmentMockRepo) Update(appointment appointment.Appointment) (*string, error) {
	var id string
	id = appointment.ID
	return &id, nil
}

// FindByID
func (r AppointmentMockRepo) FindByID(id string) (*appointment.Appointment, error) {
	appointmt := appointment.Appointment{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Now(),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Now(),
	}
	if id == appointmt.ID {
		return &appointmt, nil
	} else {
		return nil, appointment.ErrNoRecord
	}
}

// FindByPatientID
func (r AppointmentMockRepo) FindByPatientID(patientID string) ([]appointment.Appointment, error) {
	appointmt := appointment.Appointment{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Now(),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Now(),
	}
	if patientID == appointmt.PatientID {
		appointmts := []appointment.Appointment{
			appointmt,
		}
		return appointmts, nil
	} else {
		return nil, appointment.ErrNoRecord
	}
}

// FindByDoctorID
func (r AppointmentMockRepo) FindByDoctorID(doctorID string) ([]appointment.Appointment, error) {
	appointmt := appointment.Appointment{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Now(),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Now(),
	}
	if doctorID == appointmt.DoctorID {
		appointmts := []appointment.Appointment{
			appointmt,
		}
		return appointmts, nil
	} else {
		return nil, appointment.ErrNoRecord
	}
}

// FindByDate
func (r AppointmentMockRepo) FindByDate(dateTime time.Time) ([]appointment.Appointment, error) {
	appointmt := appointment.Appointment{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}
	searchY, searchM, searchD := dateTime.Date()
	appointmtY, appointmtM, appointmtD := appointmt.DateTime.Date()
	if (searchY == appointmtY) && (searchM == appointmtM) && (searchD == appointmtD) {
		appointmts := []appointment.Appointment{
			appointmt,
		}
		return appointmts, nil
	} else {
		return nil, appointment.ErrNoRecord
	}
}

// FindByMonth
func (r AppointmentMockRepo) FindByInterval(s,
	e time.time) ([]appointment.Appointment, error) {
	var db []appointment.Appointment
	db = append(db, appointment.Appointment{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
		DateTime:    time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "7fef3c47-a01a-42a6-ac45-27a440596751",
		DateTime:    time.Date(2020, 9, 8, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "19f66dc6-b5c8-497b-bba2-b982bb85ded8",
		DateTime:    time.Date(2020, 9, 9, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 10, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "723e2fa0-70a9-4c20-89d9-b5f69405b772",
		DateTime:    time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "5792340a-8c35-4183-a388-2459a8e0295a",
		DateTime:    time.Date(2020, 9, 11, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	},
	)
	start := s.Sub(1 * time.Minute)
	end := e.Add(1 * time.Minute)
	var res []appointment.Appointment
	for _, el := range db {
		if el.DateTime.After(start) && el.DateTime.Before(end) {
			res := append(res, el)
		}
	}
	return []appointment.Appointment{}, nil
}

// GetAll
func (r AppointmentMockRepo) GetAll(cursor string, after bool,
	pgSize int) ([]appointment.Appointment, bool, error) {
	var db []appointment.Appointment
	db = append(db, appointment.Appointment{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
		DateTime:    time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "7fef3c47-a01a-42a6-ac45-27a440596751",
		DateTime:    time.Date(2020, 9, 8, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "19f66dc6-b5c8-497b-bba2-b982bb85ded8",
		DateTime:    time.Date(2020, 9, 9, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 10, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "723e2fa0-70a9-4c20-89d9-b5f69405b772",
		DateTime:    time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, appointment.Appointment{
		ID:          "5792340a-8c35-4183-a388-2459a8e0295a",
		DateTime:    time.Date(2020, 9, 11, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	},
	)
	var res []appointment.Appointment
	var hasMore bool
	hasMore = false
	var respSize int
	if cursor == "" {
		if pgSize >= len(db) {
			respSize = len(db)
			hasMore = false
		} else {
			respSize = pgSize
			hasMore = true
		}
		for i := 0; i < respSize; i++ {
			res = append(res, db[i])
		}
		return res, hasMore, nil
	}
	pos := r.findPos(db, cursor)
	if pos == -1 {
		return nil, false, appointment.ErrNoRecord
	}

	if after {
		start := pos + 1
		for i := start; i < (start + pgSize); i++ {
			res = append(res, db[i])
		}
		if (len(db) - pos) > pgSize {
			hasMore = true
		}
	} else {
		start := pos - pgSize
		if start < 0 {
			start = 0
		}
		for i := start; i <= (pos - 1); i++ {
			res = append(res, db[i])
		}
		if pos > pgSize {
			hasMore = true
		}
	}
	return res, hasMore, nil
}
func (r AppointmentMockRepo) findPos(appointmts []appointment.Appointment, id string) int {
	for i, el := range appointmts {
		if el.ID == id {
			return i
		}
	}
	return -1
}
