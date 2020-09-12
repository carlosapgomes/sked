package services_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/patient"
	"carlosapgomes.com/sked/internal/services"
	uuid "github.com/satori/go.uuid"
)

func TestPatientCreate(t *testing.T) {
	repo := mocks.NewPatientRepo()
	svc := services.NewPatientService(repo)

	tests := []struct {
		name      string
		address   string
		city      string
		state     string
		phones    []string
		createdBy string
		wantError []byte
	}{
		{"Valid patient", "Valid Street, 23", "Main City", "ST", []string{"12345"}, "7f064a4e-d3bd-48a6-a305-accf4743a94f", nil},
		{"Bad uuid", "Valid Street, 22", "Main City", "ST", []string{"12345"}, "7f064a4e-d3bd-48a6-a305-accf4743a94f", []byte("repository ID not equal to new patient ID")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.Create(tt.name, tt.address, tt.city, tt.state, tt.phones, tt.createdBy)

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

func TestFindPatientByID(t *testing.T) {
	repo := mocks.NewPatientRepo()
	svc := services.NewPatientService(repo)

	testCases := []struct {
		desc        string
		patientID   string
		wantPatient *patient.Patient
		wantError   error
	}{
		{
			desc:      "Valid Patient",
			patientID: "dcce1beb-aee6-4a4d-b724-94d470817323",
			wantPatient: &patient.Patient{
				ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
				Name:      "Alice Jones",
				Phones:    []string{"6544332135"},
				CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantError: nil,
		},
		{
			desc:        "Non-existing ID",
			patientID:   "d1700797-42d4-4fe4-8fc2-60cda46f2448",
			wantPatient: nil,
			wantError:   patient.ErrNoRecord,
		},
		{
			desc:        "Invalid ID",
			patientID:   "d1700797-42d460cda46f2448",
			wantPatient: nil,
			wantError:   patient.ErrInvalidInputSyntax,
		},
		{
			desc:        "empty ID",
			patientID:   "",
			wantPatient: nil,
			wantError:   patient.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			u, err := svc.FindByID(tC.patientID)

			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if !reflect.DeepEqual(u, tC.wantPatient) {
				t.Errorf("want \n%v\n; got \n%v\n", tC.wantPatient, u)
			}

		})
	}
}

func TestUpdatePatientName(t *testing.T) {
	repo := mocks.NewPatientRepo()
	svc := services.NewPatientService(repo)

	testCases := []struct {
		desc           string
		id             string
		newPatientName string
		wantError      error
	}{
		{
			desc:           "Valid Patient",
			id:             "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newPatientName: "Johnny Silva",
			wantError:      nil,
		},
		{
			desc:           "Empty new patient name",
			id:             "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newPatientName: "",
			wantError:      patient.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.UpdateName(tC.id, tC.newPatientName)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				patient, err := svc.FindByID(tC.id)
				if err != nil {
					t.Error("Could not find patient")
				}
				if (patient != nil) && (patient.Name != tC.newPatientName) {
					t.Errorf("want \n%v\n; got \n%v\n", tC.newPatientName, patient.Name)
				}
			}
		})
	}
}

func TestPatientGetAll(t *testing.T) {
	repo := mocks.NewPatientRepo()
	svc := services.NewPatientService(repo)
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
			wantContainID: "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		},
		{
			desc:          "Valid Cursor After",
			before:        "NjhiMWQ1ZTItMzlkZC00NzEzLTg2MzEtYTA4MTAwMzgzYTBm",
			after:         "",
			pgSize:        2,
			wantSize:      1,
			hasMore:       false,
			wantError:     nil,
			wantContainID: "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		},
		{
			desc:          "Valid Cursor Before",
			before:        "",
			after:         "NjhiMWQ1ZTItMzlkZC00NzEzLTg2MzEtYTA4MTAwMzgzYTBm",
			pgSize:        2,
			wantSize:      2,
			hasMore:       true,
			wantError:     nil,
			wantContainID: "dcce1beb-aee6-4a4d-b724-94d470817323",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			cursor, err := svc.GetAll(tC.before, tC.after, tC.pgSize)
			if err != tC.wantError {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if cursor != nil && len(cursor.Patients) != tC.wantSize {
				t.Errorf("Want %v; got %v\n", tC.wantSize, len(cursor.Patients))
			}
			if tC.hasMore && !(cursor.HasAfter || cursor.HasBefore) {
				t.Errorf("want %v; got %v\n", tC.hasMore, (cursor.HasAfter || cursor.HasBefore))
			}
			var contain bool
			for _, p := range cursor.Patients {
				if p.ID == tC.wantContainID {
					contain = true
				}
			}
			if !contain {
				t.Errorf("Want response to contain %v ID;  but it did not\n", tC.wantContainID)
			}
		})
	}
}

func TestPatientFindByName(t *testing.T) {
	repo := mocks.NewPatientRepo()
	svc := services.NewPatientService(repo)
	var tests = []struct {
		name           string
		nameToSearch   string
		wantResNumber  int
		wantResContain string
		wantError      error
	}{
		{
			"Valid name",
			"Tim",
			1,
			"Tim Berners-Lee",
			nil,
		},
		{
			"Unknown name",
			"John",
			0,
			"",
			patient.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.FindByName(tt.nameToSearch)
			if err != nil && err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (res != nil) && (len(*res) != tt.wantResNumber) {
				t.Errorf("want %d results but got %d", tt.wantResNumber, len(*res))
			}
			if (res != nil) && (len(*res) > 0) && (len(tt.wantResContain) > 0) {
				var contains bool
				contains = false
				for _, u := range *res {
					if strings.Contains(strings.ToLower(u.Name), strings.ToLower(tt.wantResContain)) {
						contains = true
						t.Log("response contains desired result")
					}
				}
				if !contains {
					t.Errorf("Want results contains %v; but got nothing similar", tt.wantResContain)
				}
			}

		})
	}
}
