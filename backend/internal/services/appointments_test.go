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
		wantError   []byte
	}{
		{"Valid appointment", "John Doe", "22070f56-5d52-43f0-9f59-5de61c1db506", "Dr House", "f06244b9-97e5-4f1a-bae0-3b6da7a0b604", "some notes", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.Create(tt.dateTime, tt.patientName, tt.patientID, tt.doctorName, tt.doctorID, tt.notes)

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
