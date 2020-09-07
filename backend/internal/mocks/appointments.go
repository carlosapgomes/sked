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

//NewUserRepo returns a mocked repository
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
func (r AppointmentMockRepo) FindByPatientID(patientID string) ([]*appointment.Appointment, error) {
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
		appointmts := []*appointment.Appointment{
			&appointmt,
		}
		return appointmts, nil
	} else {
		return nil, appointment.ErrNoRecord
	}
}

// FindByDoctorID
func (r AppointmentMockRepo) FindByDoctorID(doctorID string) ([]*appointment.Appointment, error) {
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
		appointmts := []*appointment.Appointment{
			&appointmt,
		}
		return appointmts, nil
	} else {
		return nil, appointment.ErrNoRecord
	}
}

// FindByDate
func (r AppointmentMockRepo) FindByDate(dateTime time.Time) ([]*appointment.Appointment, error) {
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
	if dateTime.Equal(appointmt.DateTime) {
		appointmts := []*appointment.Appointment{
			&appointmt,
		}
		return appointmts, nil
	} else {
		return nil, appointment.ErrNoRecord
	}
}

// GetAll
func (r AppointmentMockRepo) GetAll(cursor string, after bool, pgSize int) (*[]appointment.Appointment, bool, error) {
	var res []appointment.Appointment
	res = append(res, appointment.Appointment{
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
	return &res, false, nil
}
