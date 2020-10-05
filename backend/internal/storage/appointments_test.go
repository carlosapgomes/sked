package storage_test

import (
	"errors"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
	"carlosapgomes.com/sked/internal/storage"
)

func TestCreateAppointment(t *testing.T) {
	var tests = []struct {
		name         string
		newAppointmt *appointment.Appointment
		wantError    error
	}{
		{"Valid Appointment",
			&appointment.Appointment{
				ID:          "60fa2009-e492-459d-bace-fad9da6831bf",
				DateTime:    time.Now(),
				PatientName: "John Doe",
				PatientID:   "c753a381-7642-4709-876f-57b16a5c6a6c",
				DoctorName:  "Dr House",
				DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:       "some notes",
				Canceled:    false,
				Completed:   false,
				CreatedBy:   "896d45e7-b544-41da-aa3f-f59a321fcdb9",
				CreatedAt:   time.Now(),
				UpdatedBy:   "896d45e7-b544-41da-aa3f-f59a321fcdb9",
				UpdatedAt:   time.Now(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			id, err := repo.Create(*tt.newAppointmt)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (id == nil) && (tt.wantError == nil) {
				t.Errorf("want %s; got id = nil", tt.newAppointmt.ID)
			}
			if tt.wantError == nil &&
				id != nil &&
				*id != tt.newAppointmt.ID {
				t.Errorf("want \n%v\n; got \n%v\n",
					tt.newAppointmt.ID, id)
			}
		})
	}
}

func TestUpdateAppointment(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt appointment.Appointment
		wantError error
	}{
		{"Valid Update",
			appointment.Appointment{
				ID:          "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:    time.Date(2020, 9, 8, 12, 0, 0, 0, time.UTC),
				PatientName: "John Doe",
				PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:  "Dr House",
				DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:       "some new notes",
				Canceled:    false,
				Completed:   true,
				CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:   time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
			},
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			id, err := repo.Update(tt.appointmt)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (id == nil) && (tt.wantError == nil) {
				t.Errorf("want %s; got id = nil", tt.appointmt.ID)
			}

		})
	}
}

func TestFindAppointmentByID(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt appointment.Appointment
		wantError error
	}{
		{"Valid Update",
			appointment.Appointment{
				ID:          "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:    time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName: "John Doe",
				PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:  "Dr House",
				DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:       "some notes",
				Canceled:    false,
				Completed:   false,
				CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
			},
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			appointmt, err := repo.FindByID(tt.appointmt.ID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if *appointmt != tt.appointmt {
				t.Errorf("want %s; got id = nil", tt.appointmt.ID)
			}

		})
	}

}

func TestFindAppointmentsByPatientID(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt appointment.Appointment
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			appointment.Appointment{
				ID:          "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:    time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName: "John Doe",
				PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:  "Dr House",
				DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:       "some notes",
				Canceled:    false,
				Completed:   false,
				CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
			},
			5,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			appointmts, err := repo.FindByPatientID(tt.appointmt.PatientID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(appointmts) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(appointmts))
			}
		})
	}
}

func TestFindAppointmentsByDoctorID(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt appointment.Appointment
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			appointment.Appointment{
				ID:          "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:    time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName: "John Doe",
				PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:  "Dr House",
				DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:       "some notes",
				Canceled:    false,
				Completed:   false,
				CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
			},
			5,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			appointmts, err := repo.FindByDoctorID(tt.appointmt.DoctorID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(appointmts) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(appointmts))
			}
		})
	}
}

func TestFindAppointmentsByDate(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt appointment.Appointment
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			appointment.Appointment{
				ID:          "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:    time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName: "John Doe",
				PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:  "Dr House",
				DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:       "some notes",
				Canceled:    false,
				Completed:   false,
				CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
			},
			1,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			appointmts, err := repo.FindByDate(tt.appointmt.DateTime)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(appointmts) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(appointmts))
			}
		})
	}

}
