package storage_test

import (
	"errors"
	"reflect"
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
			name: "Valid User",
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
		//{
		//name: "Bad userID",
		//newUser: &patient.Patient{
		//ID:        "5b28f739-e372-4622-9992f3c7b0e9",
		//Name:      "Muhamed Ali",
		//Address:   "Main Street, 45",
		//HashedPw:  []byte("$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG"),
		//CreatedAt: time.Now().UTC(),
		//UpdatedAt: time.Now().UTC(),
		//Active:    true,
		//Roles:     []string{},
		//},
		//wantError: user.ErrInvalidInputSyntax,
		//},
		//{
		//name: "Duplicate ID",
		//newUser: &patient.Patient{
		//ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
		//Name:      "Alice Jones",
		//Address:   "alice@example.com",
		//HashedPw:  []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
		//CreatedAt: time.Now().UTC(),
		//UpdatedAt: time.Now().UTC(),
		//Active:    true,
		//Roles:     []string{},
		//},
		//wantError: user.ErrDuplicateField,
		//},
		//{
		//name: "Duplicate email",
		//newUser: &patient.Patient{
		//ID:        "75499ef5-bde6-4f39-81b8-daf181942741",
		//Name:      "Alice Jones",
		//Address:   "alice@example.com",
		//HashedPw:  []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
		//CreatedAt: time.Now().UTC(),
		//UpdatedAt: time.Now().UTC(),
		//Active:    true,
		//Roles:     []string{},
		//},
		//wantError: user.ErrDuplicateField,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgPatientRepository(db)

			uid, err := repo.Create(*tt.newPatient)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (uid == nil) && (tt.wantError == nil) {
				t.Errorf("want %s; got uid = nil", tt.newPatient.ID)
			}
			//if uid != nil {
			//user, _ := repo.FindByID(*uid)
			//if (tt.newPatient.Name != user.Name) || (tt.newPatient.Email != user.Email) {
			//t.Errorf("want \n%v\n; got \n%v\n", tt.newPatient.Name, user.Name)
			//}
			//}
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

			repo := storage.NewPgUserRepository(db)

			user, err := repo.FindByID(tt.patientID)

			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if !reflect.DeepEqual(user, tt.wantPatient) {
				t.Errorf("want \n%v\n; got \n%v\n", tt.wantPatient, user)
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
			desc:      "Invalid uid",
			id:        "68b1d5e2-8631-a08100383a0f",
			newName:   "Johnny Smith",
			wantError: errors.New("Any error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)
			err := repo.UpdateName(tC.id, tC.newName)
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
		desc             string
		cursor           string
		after            bool
		pgSize           int
		wantSize         int
		hasMore          bool
		wantError        error
		wantContainEmail string
	}{
		{
			desc:             "Valid Page",
			cursor:           "",
			after:            true,
			pgSize:           6,
			wantSize:         6,
			hasMore:          false,
			wantError:        nil,
			wantContainEmail: "spongebob@somewhere.com",
		},
		{
			desc:             "Valid Cursor After",
			cursor:           "bobama@somewhere.com",
			after:            true,
			pgSize:           2,
			wantSize:         2,
			hasMore:          true,
			wantError:        nil,
			wantContainEmail: "spongebob@somewhere.com",
		},
		{
			desc:             "Valid Cursor Before",
			cursor:           "bobama@somewhere.com",
			after:            false,
			pgSize:           2,
			wantSize:         2,
			hasMore:          false,
			wantError:        nil,
			wantContainEmail: "alice@example.com",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgUserRepository(db)

			users, hasMore, err := repo.GetAll(tC.cursor, tC.after, tC.pgSize)
			if err != tC.wantError {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if users != nil && len(*users) != tC.wantSize {
				t.Errorf("Want %v; got %v\n", tC.wantSize, len(*users))
			}
			if tC.hasMore != hasMore {
				t.Errorf("Want %v; got %v\n", tC.hasMore, hasMore)
			}
			var contain bool
			for _, u := range *users {
				t.Logf("%v\n", u.Email)
				if u.Email == tC.wantContainEmail {
					contain = true
				}
			}
			if !contain {
				t.Errorf("Want response to contain %v address;  but it did not\n", tC.wantContainEmail)
			}
		})
	}
}