package storage_test

import (
	"errors"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/storage"
	"carlosapgomes.com/sked/internal/surgery"
)

func TestCreateSurgery(t *testing.T) {
	var tests = []struct {
		name         string
		newAppointmt *surgery.Surgery
		wantError    error
	}{
		{"Valid Surgery",
			&surgery.Surgery{
				ID:              "60fa2009-e492-459d-bace-fad9da6831bf",
				DateTime:        time.Now(),
				PatientName:     "John Doe",
				PatientID:       "c753a381-7642-4709-876f-57b16a5c6a6c",
				DoctorName:      "Dr House",
				DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:           "some notes",
				ProposedSurgery: "saphenectomy",
				Canceled:        false,
				Done:            false,
				CreatedBy:       "896d45e7-b544-41da-aa3f-f59a321fcdb9",
				CreatedAt:       time.Now(),
				UpdatedBy:       "896d45e7-b544-41da-aa3f-f59a321fcdb9",
				UpdatedAt:       time.Now(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
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

func TestUpdateSurgery(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt surgery.Surgery
		wantError error
	}{
		{"Valid Update",
			surgery.Surgery{
				ID:              "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:        time.Date(2020, 9, 8, 12, 0, 0, 0, time.UTC),
				PatientName:     "John Doe",
				PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:      "Dr House",
				DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:           "some new notes",
				ProposedSurgery: "saphenectomy",
				Canceled:        false,
				Done:            true,
				CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:       time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
			},
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
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

func TestFindSurgeryByID(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt surgery.Surgery
		wantError error
	}{
		{"Valid Update",
			surgery.Surgery{
				ID:              "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:        time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName:     "John Doe",
				PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:      "Dr House",
				DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:           "some notes",
				ProposedSurgery: "saphenectomy",
				Canceled:        false,
				Done:            false,
				CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
			},
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
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

func TestFindSurgerysByPatientID(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt surgery.Surgery
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			surgery.Surgery{
				ID:              "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:        time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName:     "John Doe",
				PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:      "Dr House",
				DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:           "some notes",
				ProposedSurgery: "saphenectomy",
				Canceled:        false,
				Done:            false,
				CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
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
			repo := storage.NewPgSurgeryRepository(db)
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

func TestFindSurgerysByDoctorID(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt surgery.Surgery
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			surgery.Surgery{
				ID:              "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:        time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName:     "John Doe",
				PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:      "Dr House",
				DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:           "some notes",
				ProposedSurgery: "saphenectomy",
				Canceled:        false,
				Done:            false,
				CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
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
			repo := storage.NewPgSurgeryRepository(db)
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

func TestFindSurgerysByDate(t *testing.T) {
	var tests = []struct {
		name      string
		appointmt surgery.Surgery
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			surgery.Surgery{
				ID:              "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
				DateTime:        time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
				PatientName:     "John Doe",
				PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
				DoctorName:      "Dr House",
				DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:           "some notes",
				ProposedSurgery: "saphenectomy",
				Canceled:        false,
				Done:            false,
				CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
				UpdatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
				UpdatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
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
			repo := storage.NewPgSurgeryRepository(db)
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

func TestGetAllSurgerys(t *testing.T) {
	var tests = []struct {
		name        string
		cursorID    string
		next        bool
		pgSize      int
		wantSize    int
		wantHasMore bool
		wantError   error
	}{
		{
			"Valid Update",
			"723e2fa0-70a9-4c20-89d9-b5f69405b772",
			true,
			3,
			3,
			true,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
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
		})
	}
}
