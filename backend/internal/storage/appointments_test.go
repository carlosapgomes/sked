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
				ID:          "f80bcdf7-6a95-4595-a9ce-f411b969ab51",
				DateTime:    time.Now(),
				PatientName: "John Doe",
				PatientID:   "c753a381-7642-4709-876f-57b16a5c6a6c",
				DoctorName:  "Dr House",
				DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
				Notes:       "some notes",
				CreatedBy:   "896d45e7-b544-41da-aa3f-f59a321fcdb9",
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
			if id != nil {
				a, _ := repo.FindByID(*id)
				if (tt.newAppointmt.PatientName != a.PatientName) ||
					(tt.newAppointmt.DoctorName != a.DoctorName) {
					t.Errorf("want \n%v\n; got \n%v\n",
						tt.newAppointmt.PatientName, a.PatientName)
					t.Errorf("want \n%v\n; got \n%v\n",
						tt.newAppointmt.DoctorName, a.DoctorName)
				}
			}
		})
	}
}
