package services_test

import (
	"bytes"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/surgery"
	uuid "github.com/satori/go.uuid"
)

func TestSurgeryCreate(t *testing.T) {
	repo := mocks.NewSurgeryRepo()
	userSvc := services.NewUserService(mocks.NewUserRepo())
	svc := services.NewSurgeryService(repo, userSvc)

	tests := []struct {
		name            string
		dateTime        time.Time
		patientName     string
		patientID       string
		doctorName      string
		doctorID        string
		notes           string
		proposedSurgery string
		createdByID     string
		wantError       []byte
	}{
		{"Valid surgery", time.Now(), "John Doe",
			"22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House",
			"f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes",
			"saphenectomy", "ecadbb28-14e6-4560-8574-809c6c54b9cb", nil},
		{"Invalid patientID format", time.Now(), "John Doe",
			"22070f56--43f0-9f59-5de61c1db506", "Dr House",
			"f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes",
			"saphenectomy", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
			[]byte("invalid input syntax")},
		{"Invalid doctorID format", time.Now(), "John Doe",
			"22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House",
			"f06244b9-97e5--bae0-3b6da7a0b604", "some notes", "saphenectomy",
			"10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
			[]byte("invalid input syntax")},
		{"Invalid createdByID format", time.Now(), "John Doe",
			"22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House",
			"f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes",
			"saphenectomy", "10b9ad06-e86d--acb1-d7e268d1f21a",
			[]byte("invalid input syntax")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.Create(tt.dateTime, tt.patientName, tt.patientID, tt.doctorName, tt.doctorID, tt.notes, tt.proposedSurgery, tt.createdByID)

			if (tt.wantError != nil) && (err != nil) {
				//t.Log("wantError and error != nil")
				e := err.Error()
				if !bytes.Contains([]byte(e), tt.wantError) {
					t.Errorf("want error msg %s to contain %q", e, tt.wantError)
				}
			}
			if (tt.wantError == nil) && (err != nil) {
				//t.Log("wantError == nil and error != nil")
				t.Errorf("want error equal 'nil'; got %s", err)
			}
			if (tt.wantError == nil) && (err == nil) {
				//t.Log("wantError and error == nil")
				if id != nil {
					_, err = uuid.FromString(*id)
					if err != nil {
						t.Errorf("want a valid uuid; got %s", *id)
					}
				} else {
					t.Errorf("want id not equal to 'nil'")
				}
			}
		})
	}

}

func TestSurgeryUpdate(t *testing.T) {
	repo := mocks.NewSurgeryRepo()
	userSvc := services.NewUserService(mocks.NewUserRepo())
	svc := services.NewSurgeryService(repo, userSvc)

	tests := []struct {
		name            string
		id              string
		dateTime        time.Time
		patientName     string
		patientID       string
		doctorName      string
		doctorID        string
		notes           string
		proposedSurgery string
		canceled        bool
		done            bool
		updatedBy       string
		wantError       []byte
	}{
		{"Valid surgery",
			"e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
			time.Now(),
			"John Doe",
			"22070f56-5d52-43f0-9f59-5de61c1db506",
			"Dr House",
			"f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
			"some notes",
			"saphenectomy",
			false,
			false,
			"10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
			nil},
		{"Invalid updatedBy",
			"e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
			time.Now(),
			"John Doe",
			"22070f56-5d52-43f0-9f59-5de61c1db506",
			"Dr House",
			"f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
			"some notes",
			"saphenectomy",
			false,
			false,
			"10b9ad06-4a85-acb1",
			[]byte("invalid input syntax")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			surgery := surgery.Surgery{
				ID:              tt.id,
				DateTime:        tt.dateTime,
				PatientName:     tt.patientName,
				PatientID:       tt.patientID,
				DoctorName:      tt.doctorName,
				DoctorID:        tt.doctorID,
				Notes:           tt.notes,
				ProposedSurgery: tt.proposedSurgery,
				Canceled:        tt.canceled,
				Done:            tt.done,
				UpdatedBy:       tt.updatedBy,
			}
			id, err := svc.Update(surgery)
			if tt.wantError != nil {
				if err != nil {
					t.Log("wantError and error != nil")
					e := err.Error()
					if !bytes.Contains([]byte(e), tt.wantError) {
						t.Errorf("want error msg %s to contain %q", e, tt.wantError)
					}
				} else {
					t.Errorf("want error msg nil to contain %q", tt.wantError)
				}

			}
			if id != nil {
				if *id != surgery.ID {
					t.Errorf("want id %s but received %s", tt.id, *id)
				}
			} else if tt.wantError == nil {
				t.Errorf("received id is nil")
			}
		})
	}
}

func TestSurgeryFindByID(t *testing.T) {
	repo := mocks.NewSurgeryRepo()
	userSvc := services.NewUserService(mocks.NewUserRepo())
	svc := services.NewSurgeryService(repo, userSvc)

	tests := []struct {
		name        string
		ID          string
		patientName string
		patientID   string
		doctorName  string
		doctorID    string
		notes       string
		createdByID string
		wantError   []byte
	}{
		{"Valid surgeryID", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", nil},
		{"Invalid surgeryID", "e521798b-9f33-9677ed1cd1ae", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("invalid input syntax")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			surgery, err := svc.FindByID(tt.ID)
			if tt.wantError != nil {
				if err != nil {
					t.Log("wantError and error != nil")
					e := err.Error()
					if !bytes.Contains([]byte(e), tt.wantError) {
						t.Errorf("want error msg %s to contain %q", e, tt.wantError)
					}
				} else {
					t.Errorf("want error msg nil to contain %q", tt.wantError)
				}
			}
			if surgery != nil {
				if surgery.ID != tt.ID {
					t.Errorf("want id %s but received %s", tt.ID, surgery.ID)
				}
				if surgery.PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", surgery.PatientID, tt.patientID)
				}
				if surgery.DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", surgery.DoctorID, tt.doctorID)
				}

			}
		})
	}
}

func TestSurgeryFindByPatientID(t *testing.T) {
	repo := mocks.NewSurgeryRepo()
	userSvc := services.NewUserService(mocks.NewUserRepo())
	svc := services.NewSurgeryService(repo, userSvc)

	tests := []struct {
		name        string
		ID          string
		patientName string
		patientID   string
		doctorName  string
		doctorID    string
		notes       string
		createdByID string
		wantError   []byte
	}{
		{"Valid patientID", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", nil},
		{"Invalid patientID", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", "John Doe", "22070f56-5d52-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("invalid input syntax")},
		{"Valid patientID without surgery", "", "John Doe", "c49a4ead-73de-46d9-92c6-9418043ae0d8", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("no matching record found")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			surgeries, err := svc.FindByPatientID(tt.patientID)
			if tt.wantError != nil {
				if err != nil {
					t.Log("wantError and error != nil")
					e := err.Error()
					if !bytes.Contains([]byte(e), tt.wantError) {
						t.Errorf("want error msg %s to contain %q", e, tt.wantError)
					}
				} else {
					t.Errorf("want error msg nil to contain %q", tt.wantError)
				}
			}
			if tt.wantError == nil && len(surgeries) == 0 {
				t.Errorf("want size of surgeries list > 0")
			}
			if surgeries != nil {
				if surgeries[0].ID != tt.ID {
					t.Errorf("want surgery id %s but received %s", tt.ID, surgeries[0].ID)
				}
				if surgeries[0].PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", surgeries[0].PatientID, tt.patientID)
				}
				if surgeries[0].DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", surgeries[0].DoctorID, tt.doctorID)
				}

			}
		})
	}
}

func TestSurgeryFindByDoctorID(t *testing.T) {
	repo := mocks.NewSurgeryRepo()
	userSvc := services.NewUserService(mocks.NewUserRepo())
	svc := services.NewSurgeryService(repo, userSvc)

	tests := []struct {
		name        string
		ID          string
		patientName string
		patientID   string
		doctorName  string
		doctorID    string
		notes       string
		createdByID string
		wantError   []byte
	}{
		{"Valid doctorID", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", nil},
		{"Invalid doctorID", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("invalid input syntax")},
		{"Valid doctorID without surgeries", "", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "1dc5b27f-4ff4-4b96-a80d-6702912cf0a0", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("no matching record found")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			surgeries, err := svc.FindByDoctorID(tt.doctorID)
			if tt.wantError != nil {
				if err != nil {
					t.Log("wantError and error != nil")
					e := err.Error()
					if !bytes.Contains([]byte(e), tt.wantError) {
						t.Errorf("want error msg %s to contain %q", e, tt.wantError)
					}
				} else {
					t.Errorf("want error msg nil to contain %q", tt.wantError)
				}
			}
			if tt.wantError == nil && len(surgeries) == 0 {
				t.Errorf("want size of surgeries list > 0")
			}
			if surgeries != nil {
				if surgeries[0].ID != tt.ID {
					t.Errorf("want surgery id %s but received %s", tt.ID, surgeries[0].ID)
				}
				if surgeries[0].PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", surgeries[0].PatientID, tt.patientID)
				}
				if surgeries[0].DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", surgeries[0].DoctorID, tt.doctorID)
				}

			}
		})
	}
}

func TestFindSurgeriesByMonthYear(t *testing.T) {
	userRepo := mocks.NewUserRepo()
	userSvc := services.NewUserService(userRepo)
	repo := mocks.NewSurgeryRepo()
	svc := services.NewSurgeryService(repo, userSvc)

	tests := []struct {
		name      string
		month     int
		year      int
		wantSize  int
		wantError error
	}{
		{"Valid Month/Year", 9, 2020, 6, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			surgs, err := svc.FindByMonthYear(tt.month, tt.year)
			if tt.wantError != err {
				t.Errorf("Want error %v, but got %v\n", tt.wantError, err)
			}
			if tt.wantSize != len(surgs) {
				t.Errorf("Want respose size %v, but got %v\n", tt.wantSize,
					len(surgs))
			}
		})
	}
}

func TestSurgeryFindByDate(t *testing.T) {
	repo := mocks.NewSurgeryRepo()
	userSvc := services.NewUserService(mocks.NewUserRepo())
	svc := services.NewSurgeryService(repo, userSvc)

	tests := []struct {
		name        string
		id          string
		dateTime    time.Time
		patientName string
		patientID   string
		doctorName  string
		doctorID    string
		notes       string
		createdByID string
		wantError   []byte
	}{
		{"Valid date", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", nil},
		{"Valid date without surgeries", "", time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "1dc5b27f-4ff4-4b96-a80d-6702912cf0a0", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("no matching record found")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			surgeries, err := svc.FindByDate(tt.dateTime)
			if tt.wantError != nil {
				if err != nil {
					t.Log("wantError and error != nil")
					e := err.Error()
					if !bytes.Contains([]byte(e), tt.wantError) {
						t.Errorf("want error msg %s to contain %q", e, tt.wantError)
					}
				} else {
					t.Errorf("want error msg nil to contain %q", tt.wantError)
				}
			}
			if tt.wantError == nil && len(surgeries) == 0 {
				t.Errorf("want size of surgeries list > 0")
			}
			if surgeries != nil {
				if surgeries[0].ID != tt.id {
					t.Errorf("want surgery id %s but received %s", tt.id, surgeries[0].ID)
				}
				if surgeries[0].PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", surgeries[0].PatientID, tt.patientID)
				}
				if surgeries[0].DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", surgeries[0].DoctorID, tt.doctorID)
				}

			}
		})
	}
}

func TestSurgeryGetAll(t *testing.T) {
	repo := mocks.NewSurgeryRepo()
	userSvc := services.NewUserService(mocks.NewUserRepo())
	svc := services.NewSurgeryService(repo, userSvc)
	testCases := []struct {
		desc          string
		before        string
		after         string
		pgSize        int
		wantSize      int
		hasMore       bool
		wantError     error
		wantContainID string
	}{
		{
			desc:          "Valid Page",
			before:        "",
			after:         "",
			pgSize:        6,
			wantSize:      6,
			hasMore:       false,
			wantError:     nil,
			wantContainID: "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
		},
		{
			desc:          "Valid Cursor Before",
			before:        "NWU2ZjdjZDEtZDhkMi00MGNkLTk3YTMtYWNhMDFhOTNiZmRl",
			after:         "",
			pgSize:        2,
			wantSize:      1,
			hasMore:       false,
			wantError:     nil,
			wantContainID: "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		},
		{
			desc:          "Valid Cursor After",
			before:        "",
			after:         "NWU2ZjdjZDEtZDhkMi00MGNkLTk3YTMtYWNhMDFhOTNiZmRl",
			pgSize:        2,
			wantSize:      2,
			hasMore:       true,
			wantError:     nil,
			wantContainID: "7fef3c47-a01a-42a6-ac45-27a440596751",
		},
		//{
		//desc:             "Valid Cursor Before",
		//cursor:           "bobama@somewhere.com",
		//after:            false,
		//pgSize:           2,
		//wantSize:         2,
		//hasMore:          false,
		//wantError:        nil,
		//wantContainEmail: "alice@example.com",
		//},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			cursor, err := svc.GetAll(tC.before, tC.after, tC.pgSize)
			if err != tC.wantError {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if cursor != nil && len(cursor.Surgeries) != tC.wantSize {
				t.Errorf("Want %v; got %v\n", tC.wantSize, len(cursor.Surgeries))
			}
			if tC.hasMore && !(cursor.HasNextPage || cursor.HasPreviousPage) {
				t.Errorf("want %v; got %v\n", tC.hasMore, (cursor.HasNextPage || cursor.HasPreviousPage))
			}
			var contain bool
			for _, u := range cursor.Surgeries {
				if u.ID == tC.wantContainID {
					contain = true
				}
			}
			if !contain {
				t.Errorf("Want response to contain %v ID;  but it did not\n", tC.wantContainID)
			}
		})
	}
}
