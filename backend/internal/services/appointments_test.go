package services_test

import (
	"bytes"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/services"
	uuid "github.com/satori/go.uuid"
)

func TestAppointmentCreate(t *testing.T) {
	repo := mocks.NewAppointmentRepo()
	svc := services.NewAppointmentService(repo)

	tests := []struct {
		name        string
		dateTime    time.Time
		patientName string
		patientID   string
		doctorName  string
		doctorID    string
		notes       string
		createdByID string
		wantError   []byte
	}{
		{"Valid appointment", time.Now(), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", nil},
		{"Invalid patientID format", time.Now(), "John Doe", "22070f56--43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("invalid input syntax")},
		{"Invalid doctorID format", time.Now(), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5--bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("invalid input syntax")},
		{"Invalid createdByID format", time.Now(), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d--acb1-d7e268d1f21a", []byte("invalid input syntax")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.Create(tt.dateTime, tt.patientName, tt.patientID, tt.doctorName, tt.doctorID, tt.notes, tt.createdByID)

			if (tt.wantError != nil) && (err != nil) {
				t.Log("wantError and error != nil")
				e := err.Error()
				if !bytes.Contains([]byte(e), tt.wantError) {
					t.Errorf("want error msg %s to contain %q", e, tt.wantError)
				}
			}
			if (tt.wantError == nil) && (err != nil) {
				t.Log("wantError == nil and error != nil")
				t.Errorf("want error equal 'nil'; got %s", err)
			}
			if (tt.wantError == nil) && (err == nil) {
				t.Log("wantError and error == nil")
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

func TestAppointmentUpdate(t *testing.T) {
	repo := mocks.NewAppointmentRepo()
	svc := services.NewAppointmentService(repo)

	tests := []struct {
		name        string
		id          string
		dateTime    time.Time
		patientName string
		patientID   string
		doctorName  string
		doctorID    string
		notes       string
		updatedBy   string
		wantError   []byte
	}{
		{"Valid appointment", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", time.Now(), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", nil},
		{"Invalid updatedBy", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", time.Now(), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-4a85-acb1", []byte("invalid input syntax")},
		//{"Invalid updatedBy", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", time.Now(), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e-acb1-d7e268d1f21a", []byte("invalid input syntax")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			appointmt := appointment.Appointment{
				ID:          tt.id,
				DateTime:    tt.dateTime,
				PatientName: tt.patientName,
				PatientID:   tt.patientID,
				DoctorName:  tt.doctorName,
				DoctorID:    tt.doctorID,
				Notes:       tt.notes,
				UpdatedBy:   tt.updatedBy,
			}
			id, err := svc.Update(appointmt)
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
				if *id != appointmt.ID {
					t.Errorf("want id %s but received %s", tt.id, *id)
				}
			} else if tt.wantError == nil {
				t.Errorf("received id is nil")
			}
		})
	}
}

func TestAppointmentFindByID(t *testing.T) {
	repo := mocks.NewAppointmentRepo()
	svc := services.NewAppointmentService(repo)

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
		{"Valid appointmentID", "e521798b-9f33-4a10-8b2a-9677ed1cd1ae", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", nil},
		{"Invalid appointmentID", "e521798b-9f33-9677ed1cd1ae", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("invalid input syntax")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			appointmt, err := svc.FindByID(tt.ID)
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
			if appointmt != nil {
				if appointmt.ID != tt.ID {
					t.Errorf("want id %s but received %s", tt.ID, appointmt.ID)
				}
				if appointmt.PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", appointmt.PatientID, tt.patientID)
				}
				if appointmt.DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", appointmt.DoctorID, tt.doctorID)
				}

			}
		})
	}
}

func TestAppointmentFindByPatientID(t *testing.T) {
	repo := mocks.NewAppointmentRepo()
	svc := services.NewAppointmentService(repo)

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
		{"Valid patientID without appointement", "", "John Doe", "c49a4ead-73de-46d9-92c6-9418043ae0d8", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("no matching record found")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			appointmts, err := svc.FindByPatientID(tt.patientID)
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
			if tt.wantError == nil && len(appointmts) == 0 {
				t.Errorf("want size of appointements list > 0")
			}
			if appointmts != nil {
				if appointmts[0].ID != tt.ID {
					t.Errorf("want appointement id %s but received %s", tt.ID, appointmts[0].ID)
				}
				if appointmts[0].PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", appointmts[0].PatientID, tt.patientID)
				}
				if appointmts[0].DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", appointmts[0].DoctorID, tt.doctorID)
				}

			}
		})
	}
}

func TestAppointmentFindByDoctorID(t *testing.T) {
	repo := mocks.NewAppointmentRepo()
	svc := services.NewAppointmentService(repo)

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
		{"Valid doctorID without appointements", "", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "1dc5b27f-4ff4-4b96-a80d-6702912cf0a0", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("no matching record found")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			appointmts, err := svc.FindByDoctorID(tt.doctorID)
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
			if tt.wantError == nil && len(appointmts) == 0 {
				t.Errorf("want size of appointements list > 0")
			}
			if appointmts != nil {
				if appointmts[0].ID != tt.ID {
					t.Errorf("want appointement id %s but received %s", tt.ID, appointmts[0].ID)
				}
				if appointmts[0].PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", appointmts[0].PatientID, tt.patientID)
				}
				if appointmts[0].DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", appointmts[0].DoctorID, tt.doctorID)
				}

			}
		})
	}
}

func TestAppointmentFindByDate(t *testing.T) {
	repo := mocks.NewAppointmentRepo()
	svc := services.NewAppointmentService(repo)

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
		{"Valid date without appointements", "", time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC), "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "1dc5b27f-4ff4-4b96-a80d-6702912cf0a0", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("no matching record found")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			appointmts, err := svc.FindByDate(tt.dateTime)
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
			if tt.wantError == nil && len(appointmts) == 0 {
				t.Errorf("want size of appointements list > 0")
			}
			if appointmts != nil {
				if appointmts[0].ID != tt.id {
					t.Errorf("want appointement id %s but received %s", tt.id, appointmts[0].ID)
				}
				if appointmts[0].PatientID != tt.patientID {
					t.Errorf("want patientID = %s but got %s \n", appointmts[0].PatientID, tt.patientID)
				}
				if appointmts[0].DoctorID != tt.doctorID {
					t.Errorf("want doctorID = %s but got %s \n", appointmts[0].DoctorID, tt.doctorID)
				}

			}
		})
	}
}