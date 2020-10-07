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
		name       string
		newSurgery *surgery.Surgery
		wantError  error
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
			id, err := repo.Create(*tt.newSurgery)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (id == nil) && (tt.wantError == nil) {
				t.Errorf("want %s; got id = nil", tt.newSurgery.ID)
			}
			if tt.wantError == nil &&
				id != nil &&
				*id != tt.newSurgery.ID {
				t.Errorf("want \n%v\n; got \n%v\n",
					tt.newSurgery.ID, id)
			}
		})
	}
}

func TestUpdateSurgery(t *testing.T) {
	var tests = []struct {
		name      string
		surg      surgery.Surgery
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
			id, err := repo.Update(tt.surg)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (id == nil) && (tt.wantError == nil) {
				t.Errorf("want %s; got id = nil", tt.surg.ID)
			}

		})
	}
}

func TestFindSurgeryByID(t *testing.T) {
	var tests = []struct {
		name      string
		surg      surgery.Surgery
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
		{"Non Existing Surgery",
			surgery.Surgery{
				ID: "3b5c10cc-ca7c-46b3-a83e-060515e7e162",
			},
			surgery.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
			surg, err := repo.FindByID(tt.surg.ID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if tt.wantError == nil && *surg != tt.surg {
				t.Errorf("want %s; got id = nil", tt.surg.ID)
			}

		})
	}

}

func TestFindSurgeryByPatientID(t *testing.T) {
	var tests = []struct {
		name      string
		patientID string
		surg      surgery.Surgery
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			"22070f56-5d52-43f0-9f59-5de61c1db506",
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
		{
			"PatientID With No Scheduled Surgery",
			"3f7573a9-26a0-44ea-932e-f83f480f964f",
			surgery.Surgery{},
			0,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
			surgs, err := repo.FindByPatientID(tt.patientID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(surgs) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(surgs))
			}
			if tt.wantSize > 0 {
				hasSurg := false
				for _, s := range surgs {
					if s == tt.surg {
						hasSurg = true
					}
				}
				if !hasSurg {
					t.Errorf("did not receive the expected surgery object")
				}
			}
		})
	}
}

func TestFindSurgeryByDoctorID(t *testing.T) {
	var tests = []struct {
		name      string
		doctorID  string
		surg      surgery.Surgery
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			"f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
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
		{
			"DoctorID With No Scheduled Surgery",
			"3f7573a9-26a0-44ea-932e-f83f480f964f",
			surgery.Surgery{},
			0,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
			surgs, err := repo.FindByDoctorID(tt.doctorID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(surgs) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(surgs))
			}
			if tt.wantSize > 0 {
				hasSurg := false
				for _, s := range surgs {
					if s == tt.surg {
						hasSurg = true
					}
				}
				if !hasSurg {
					t.Errorf("did not receive the expected surgery object")
				}
			}
		})
	}
}

func TestFindSurgeryByDate(t *testing.T) {
	var tests = []struct {
		name      string
		dateTime  time.Time
		surg      surgery.Surgery
		wantSize  int
		wantError error
	}{
		{"Valid Update",
			time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
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
		{
			"Date With No Scheduled Surgery",
			time.Date(2020, 4, 6, 12, 0, 0, 0, time.UTC),
			surgery.Surgery{},
			0,
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSurgeryRepository(db)
			surgs, err := repo.FindByDate(tt.dateTime)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(surgs) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(surgs))
			}
			if tt.wantSize > 0 {
				hasSurg := false
				for _, s := range surgs {
					if s == tt.surg {
						hasSurg = true
					}
				}
				if !hasSurg {
					t.Errorf("did not receive the expected surgery object")
				}
			}
		})
	}
}

func TestGetAllSurgeries(t *testing.T) {
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
			"Next Request",
			"723e2fa0-70a9-4c20-89d9-b5f69405b772",
			true,
			3,
			3,
			true,
			nil,
		},
		{
			"Previous Request",
			"00707378-7fd5-4cbe-92e1-ca93301eda49",
			false,
			2,
			2,
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
			surgs, hasMore, err := repo.GetAll(tt.cursorID,
				tt.next, tt.pgSize)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if len(surgs) != tt.wantSize {
				t.Errorf("want answer size of %d; got %d\n",
					tt.wantSize, len(surgs))
			}
			if hasMore != tt.wantHasMore {
				t.Errorf("want hasMore = %v; got %v\n", tt.wantHasMore, hasMore)
			}
		})
	}
}
