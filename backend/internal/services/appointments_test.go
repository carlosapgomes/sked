package services_test

import (
	"bytes"
	"testing"
	"time"

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
		//{"Invalid appointmentID", "e521798b-9f33-9677ed1cd1ae", "John Doe", "22070f56--43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", "10b9ad06-e86d-4a85-acb1-d7e268d1f21a", []byte("invalid input syntax")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			appointmt, err := svc.FindByID(tt.ID)
			if (tt.wantError != nil) && (err != nil) {
				t.Log("wantError and error != nil")
				e := err.Error()
				if !bytes.Contains([]byte(e), tt.wantError) {
					t.Errorf("want error msg %s to contain %q", e, tt.wantError)
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
