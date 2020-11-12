package storage_test

import (
	"errors"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/patient"
	"carlosapgomes.com/sked/internal/storage"
)

func TestCreatePatient(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name       string
		newPatient *patient.Patient
		wantError  error
	}{
		{
			name: "Valid Patient",
			newPatient: &patient.Patient{
				ID:        "5b28f739-e372-4622-8390-9992f3c7b0e9",
				Name:      "Muhamed Ali",
				Address:   "MainStreet, 42",
				City:      "MainCity",
				State:     "MN",
				Phones:    []string{"3453452"},
				CreatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC(),
				UpdatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				UpdatedAt: time.Now().UTC(),
			},
			wantError: nil,
		},
		{
			name: "Bad patientID",
			newPatient: &patient.Patient{
				ID:        "5b28f739-e372-4622-9992f3c7b0e9",
				Name:      "Alice Jones",
				Address:   "MainStreet, 42",
				City:      "MainCity",
				State:     "MN",
				Phones:    []string{"2323234"},
				CreatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC(),
				UpdatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				UpdatedAt: time.Now().UTC(),
			},
			wantError: patient.ErrInvalidInputSyntax,
		},
		{
			name: "Duplicate ID",
			newPatient: &patient.Patient{
				ID:        "68b1d5e2-39dd-4713-8631-a08100383a0f",
				Name:      "Alice Jones",
				Address:   "MainStreet, 42",
				City:      "MainCity",
				State:     "MN",
				Phones:    []string{"2323234"},
				CreatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC(),
				UpdatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				UpdatedAt: time.Now().UTC(),
			},
			wantError: patient.ErrDuplicateField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgPatientRepository(db)

			id, err := repo.Create(*tt.newPatient)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (id == nil) && (tt.wantError == nil) {
				t.Errorf("want %s; got uid = nil", tt.newPatient.ID)
			}
		})
	}

}
func TestUpdatePatient(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name           string
		updatedPatient *patient.Patient
		wantError      error
	}{
		{
			name: "Valid Patient",
			updatedPatient: &patient.Patient{
				ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
				Name:      "Alice Mary Jones",
				Address:   "MainStreet, 44",
				City:      "MainCity",
				State:     "MN",
				Phones:    []string{"46456422"},
				CreatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC(),
				UpdatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				UpdatedAt: time.Now().UTC(),
			},
			wantError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgPatientRepository(db)

			err := repo.UpdatePatient(tt.updatedPatient)
			if err != tt.wantError {
				t.Errorf("Wanted error %v, but got %v\n",
					tt.wantError, err)
			}
			if tt.wantError == nil {
				p, err := repo.FindByID(tt.updatedPatient.ID)
				if err != nil {
					t.Error("Could not retrieve updated patient")
				}
				if (p.Name != tt.updatedPatient.Name) ||
					(p.Address != tt.updatedPatient.Address) ||
					(p.Phones[0] != tt.updatedPatient.Phones[0]) {
					t.Error("Could not update patient")
				}
			}

		})
	}

}
func TestFindPatientByName(t *testing.T) {
	tests := []struct {
		name            string
		patientName     string
		wantContainName string
		wantError       error
	}{
		{
			"Valid Single Result",
			"Alice",
			"Alice Jones",
			nil,
		},
		{
			"Valid Multiple Results",
			"Bob",
			"SpongeBob Squarepants",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgPatientRepository(db)

			patients, err := repo.FindByName(tt.patientName)

			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if patients != nil {
				contain := false
				for _, u := range *patients {
					if u.Name == tt.wantContainName {
						contain = true
					}
				}
				if !contain {
					t.Errorf("want result to contain %v, but it did not\n",
						tt.wantContainName)
				}
			} else {
				t.Logf("patients = nil\n")
			}
		})
	}

}
func TestFindPatientByID(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name        string
		patientID   string
		wantPatient *patient.Patient
		wantError   error
	}{
		{
			name:      "Valid ID",
			patientID: "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			wantPatient: &patient.Patient{
				ID:        "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fc",
				Name:      "Valid Patient",
				Address:   "Somewhere Street 42",
				City:      "Main City",
				State:     "ST",
				Phones:    []string{"6544332135"},
				CreatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Date(2019, 06, 23, 17, 25, 00, 0, time.UTC),
				UpdatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
				UpdatedAt: time.Date(2019, 06, 23, 17, 25, 00, 0, time.UTC),
			},
			wantError: nil,
		},
		{
			name:        "Non-existing ID",
			patientID:   "d1700797-42d4-4fe4-8fc2-60cda46f2448",
			wantPatient: nil,
			wantError:   patient.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgPatientRepository(db)

			p, err := repo.FindByID(tt.patientID)

			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if !(p == nil) &&
				(p.ID == tt.wantPatient.ID) &&
				(p.Name == tt.wantPatient.Name) &&
				(p.Address == tt.wantPatient.Address) &&
				(p.City == tt.wantPatient.City) &&
				(p.State == tt.wantPatient.State) &&
				(p.Phones[0] == tt.wantPatient.Phones[0]) &&
				(p.CreatedBy == tt.wantPatient.CreatedBy) &&
				(p.CreatedAt == tt.wantPatient.CreatedAt) &&
				(p.UpdatedBy == tt.wantPatient.UpdatedBy) &&
				(p.UpdatedAt == tt.wantPatient.UpdatedAt) {
				t.Errorf("want \n%v\n; got \n%v\n", tt.wantPatient, p)
			}
		})
	}
}

func TestUpdatePatientName(t *testing.T) {
	testCases := []struct {
		desc      string
		id        string
		newName   string
		wantError error
	}{
		{
			desc:      "Valid id",
			id:        "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			newName:   "Johnny Smith",
			wantError: nil,
		},
		{
			desc:      "Invalid id",
			id:        "68b1d5e2-8631-a08100383a0f",
			newName:   "Johnny Smith",
			wantError: errors.New("Any error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgPatientRepository(db)
			err := repo.UpdateName(tC.id, tC.newName,
				"dcce1beb-aee6-4a4d-b724-94d470817323")
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			p, err := repo.FindByID(tC.id)
			if p != nil && p.Name != tC.newName {
				t.Errorf("Want %v; got %v\n", tC.newName, p.Name)
			}

		})
	}
}

func TestUpdatePatientPhones(t *testing.T) {
	testCases := []struct {
		desc      string
		id        string
		newPhones []string
		updatedBy string
		wantError error
	}{
		{
			desc:      "Valid id",
			id:        "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			newPhones: []string{"214377669988"},
			updatedBy: "dcce1beb-aee6-4a4d-b724-94d470817323",
			wantError: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgPatientRepository(db)
			p, _ := repo.FindByID(tC.id)
			t.Log(p.Phones)
			t.Log(p.Name)
			t.Log(p.ID)
			err := repo.UpdatePhone(tC.id, tC.newPhones, tC.updatedBy)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			p, err = repo.FindByID(tC.id)
			if p != nil && p.Phones[0] != tC.newPhones[0] {
				t.Errorf("Want %v; got %v\n %v\n %v\n", tC.newPhones[0],
					p.Phones[0], p.Name, p.ID)
			}

		})
	}
}

func TestGetAllPatients(t *testing.T) {
	testCases := []struct {
		desc            string
		cursor          string
		next            bool
		pgSize          int
		wantSize        int
		wantHasMore     bool
		wantError       error
		wantContainName string
	}{
		{
			desc:            "Next Request",
			cursor:          "dcce1beb-aee6-4a4d-b724-94d470817323",
			next:            true,
			pgSize:          5,
			wantSize:        5,
			wantHasMore:     false,
			wantError:       nil,
			wantContainName: "SpongeBob Squarepants",
		},
		{
			desc:            "Next Request With HasMore",
			cursor:          "dcce1beb-aee6-4a4d-b724-94d470817323",
			next:            true,
			pgSize:          3,
			wantSize:        3,
			wantHasMore:     true,
			wantError:       nil,
			wantContainName: "SpongeBob Squarepants",
		},
		{
			desc:            "Previous Request",
			cursor:          "27f9802b-acb3-4852-bf97-c4ed4c3b3658",
			next:            false,
			pgSize:          3,
			wantSize:        3,
			wantHasMore:     true,
			wantError:       nil,
			wantContainName: "SpongeBob Squarepants",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgPatientRepository(db)

			patients, hasMore, err := repo.GetAll(tC.cursor, tC.next, tC.pgSize)
			if err != tC.wantError {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if patients != nil && len(*patients) != tC.wantSize {
				t.Errorf("Want %v; got %v\n", tC.wantSize, len(*patients))
			}
			if tC.wantHasMore != hasMore {
				t.Errorf("Want hasMore = %v; got %v\n", tC.wantHasMore, hasMore)
			}
			var contain bool
			for _, p := range *patients {
				//t.Logf("%v\n", p.Name)
				if p.Name == tC.wantContainName {
					contain = true
				}
			}
			if !contain {
				t.Errorf("Want response to contain %v;  but it did not\n",
					tC.wantContainName)
			}
		})
	}
}
