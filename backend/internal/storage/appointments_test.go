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
		name        string
		appointmtID string
		appointmt   appointment.Appointment
		wantError   error
	}{
		{"Valid Appointment",
			"5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
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
		{
			"Invalid Appointment",
			"ed06f7f9-5fc4-4cfe-ad71-3efd24bf748a",
			appointment.Appointment{},
			appointment.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			appointmt, err := repo.FindByID(tt.appointmtID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if err == nil && *appointmt != tt.appointmt {
				t.Errorf("want %s; got id = nil", tt.appointmt.ID)
			}

		})
	}

}

func TestFindAppointmentsByPatientID(t *testing.T) {
	var tests = []struct {
		name      string
		patientID string
		appointmt appointment.Appointment
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			"22070f56-5d52-43f0-9f59-5de61c1db506",
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
		{
			"Patient without appointments",
			"ed06f7f9-5fc4-4cfe-ad71-3efd24bf748a",
			appointment.Appointment{},
			0,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			appointmts, err := repo.FindByPatientID(tt.patientID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(appointmts) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(appointmts))
			}
			if tt.wantSize > 0 {
				hasAppointmt := false
				for _, s := range appointmts {
					if s == tt.appointmt {
						hasAppointmt = true
					}
				}
				if !hasAppointmt {
					t.Errorf("did not receive the expected appointment object")
				}
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
		{
			"Doctor without appointments",
			appointment.Appointment{
				DoctorID: "ed06f7f9-5fc4-4cfe-ad71-3efd24bf748a",
			},
			0,
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
			if tt.wantSize > 0 {
				hasAppointmt := false
				for _, s := range appointmts {
					if s == tt.appointmt {
						hasAppointmt = true
					}
				}
				if !hasAppointmt {
					t.Errorf("did not receive the expected appointment object")
				}
			}
		})
	}
}

func TestFindAppointmentsByInterval(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()
	repo := storage.NewPgAppointmentRepository(db)
	var tests = []struct {
		name      string
		start     time.Time
		end       time.Time
		wantSize  int
		wantError error
	}{
		{
			"Valid Search",
			time.Date(2020, 9, 1, 12, 0, 0, 0, time.UTC),
			time.Date(2020, 9, 30, 12, 0, 0, 0, time.UTC),
			5,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			appointmts, err := repo.FindByInterval(tt.start, tt.end)

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
		{
			"Date Without Appointments",
			appointment.Appointment{
				DateTime: time.Date(2020, 8, 7, 12, 0, 0, 0, time.UTC),
			},
			0,
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
			if tt.wantSize > 0 {
				hasAppointmt := false
				for _, s := range appointmts {
					if s == tt.appointmt {
						hasAppointmt = true
					}
				}
				if !hasAppointmt {
					t.Errorf("did not receive wanted surgery")
				}
			}
		})
	}
}

func TestGetAllAppointments(t *testing.T) {
	var tests = []struct {
		name            string
		cursorID        string
		next            bool
		pgSize          int
		wantSize        int
		wantHasMore     bool
		wantError       error
		wantContainDate time.Time
	}{
		{
			"Next request",
			"723e2fa0-70a9-4c20-89d9-b5f69405b772",
			true,
			3,
			3,
			true,
			nil,
			time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
		},
		{
			"Previous Request",
			"640583f6-7727-4024-8b49-00be8d195a23",
			false,
			2,
			2,
			true,
			nil,
			time.Date(2020, 9, 8, 12, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgAppointmentRepository(db)
			appointmts, hasMore, err := repo.GetAll(tt.cursorID,
				tt.next, tt.pgSize)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(appointmts) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(appointmts))
			}
			if hasMore != tt.wantHasMore {
				t.Errorf("want hasMore = %v; got %v\n", tt.wantHasMore, hasMore)
			}
			var contain bool
			for _, a := range appointmts {
				//t.Logf("%v\n", p.Name)
				if a.DateTime == tt.wantContainDate {
					contain = true
				}
			}
			if !contain {
				t.Errorf("Want response to contain %v;  but it did not\n",
					tt.wantContainDate)
			}
		})
	}
}
