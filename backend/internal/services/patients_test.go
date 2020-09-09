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
	"carlosapgomes.com/sked/internal/user"
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
		{"Bad uuid", "Valid Street, 22", "Main City", "ST", []string{"12345"}, "7f064a4e-d3bd-48a6-a305-accf4743a94f", []byte("repository ID not equal to new user ID")},
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
		desc      string
		userID    string
		wantUser  *user.User
		wantError error
	}{
		{
			desc:   "Valid User",
			userID: "dcce1beb-aee6-4a4d-b724-94d470817323",
			wantUser: &user.User{
				ID:                "dcce1beb-aee6-4a4d-b724-94d470817323",
				Name:              "Alice Jones",
				Email:             "alice@example.com",
				Phone:             "6544332135",
				HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
				CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Active:            true,
				EmailWasValidated: true,
				Roles:             []string{user.RoleCommon},
			},
			wantError: nil,
		},
		{
			desc:      "Non-existing ID",
			userID:    "d1700797-42d4-4fe4-8fc2-60cda46f2448",
			wantUser:  nil,
			wantError: patient.ErrNoRecord,
		},
		{
			desc:      "Invalid ID",
			userID:    "d1700797-42d460cda46f2448",
			wantUser:  nil,
			wantError: patient.ErrInvalidInputSyntax,
		},
		{
			desc:      "empty ID",
			userID:    "",
			wantUser:  nil,
			wantError: patient.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			u, err := svc.FindByID(tC.userID)

			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if !reflect.DeepEqual(u, tC.wantUser) {
				t.Errorf("want \n%v\n; got \n%v\n", tC.wantUser, u)
			}

		})
	}
}

func TestFindPatientByEmail(t *testing.T) {
	repo := mocks.NewPatientRepo()
	svc := services.NewPatientService(repo)

	testCases := []struct {
		desc      string
		userEmail string
		wantUser  *user.User
		wantError error
	}{
		{
			desc:      "Valid User",
			userEmail: "alice@example.com",
			wantUser: &user.User{
				ID:                "dcce1beb-aee6-4a4d-b724-94d470817323",
				Name:              "Alice Jones",
				Email:             "alice@example.com",
				Phone:             "6544332135",
				HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
				CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Active:            true,
				EmailWasValidated: true,
				Roles:             []string{user.RoleCommon},
			},
			wantError: nil,
		},
		{
			desc:      "Non-existing user",
			userEmail: "joe@nowhere.com",
			wantUser:  nil,
			wantError: patient.ErrNoRecord,
		},
		{
			desc:      "empty Email",
			userEmail: "",
			wantUser:  nil,
			wantError: patient.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			u, err := svc.FindByName(tC.userEmail)

			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if !reflect.DeepEqual(u, tC.wantUser) {
				t.Errorf("want \n%v\n; got \n%v\n", tC.wantUser, u)
			}

		})
	}
}

func TestUpdatePatientName(t *testing.T) {
	repo := mocks.NewPatientRepo()
	svc := services.NewPatientService(repo)

	testCases := []struct {
		desc        string
		uid         string
		newUserName string
		wantError   error
	}{
		{
			desc:        "Valid User",
			uid:         "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newUserName: "Johnny Silva",
			wantError:   nil,
		},
		{
			desc:        "Empty new user name",
			uid:         "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newUserName: "",
			wantError:   patient.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.UpdateName(tC.uid, tC.newUserName)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				user, err := svc.FindByID(tC.uid)
				if err != nil {
					t.Error("Could not find user")
				}
				if (user != nil) && (user.Name != tC.newUserName) {
					t.Errorf("want \n%v\n; got \n%v\n", tC.newUserName, user.Name)
				}
			}
		})
	}
}

//func TestPatientUpdatePhone(t *testing.T) {
//repo := mocks.NewPatientRepo()
//svc := services.NewPatientService(repo)

//testCases := []struct {
//desc      string
//uid       string
//newPhone  string
//wantError error
//}{
//{
//desc:      "Valid User",
//uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
//newPhone:  "3453453452",
//wantError: nil,
//},
//{
//desc:      "Empty new user name",
//uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
//newPhone:  "",
//wantError: patient.ErrInvalidInputSyntax,
//},
//}
//for _, tC := range testCases {
//t.Run(tC.desc, func(t *testing.T) {
//err := svc.UpdatePhone(tC.uid, tC.newPhone)
//if err != tC.wantError {
//t.Errorf("want %v; got %v", tC.wantError, err)
//}
//if tC.wantError == nil {
//user, err := svc.FindByID(tC.uid)
//if err != nil {
//t.Error("Could not find user")
//}
//if (user != nil) && (user.Phones != tC.newPhone) {
//t.Errorf("want \n%v\n; got \n%v\n", tC.newPhone, user.Phones)
//}
//}
//})
//}
//}

func TestPatientGetAll(t *testing.T) {
	testCases := []struct {
		desc             string
		cursor           string
		after            bool
		pgSize           int
		wantSize         int
		wantError        error
		wantContainEmail string
	}{
		{
			desc:             "Valid Page",
			cursor:           "",
			after:            true,
			pgSize:           6,
			wantSize:         6,
			wantError:        nil,
			wantContainEmail: "spongebob@somewhere.com",
		},
		{
			desc:             "Valid Cursor After",
			cursor:           "bobama@somewhere.com",
			after:            true,
			pgSize:           2,
			wantSize:         2,
			wantError:        nil,
			wantContainEmail: "spongebob@somewhere.com",
		},
		{
			desc:             "Valid Cursor Before",
			cursor:           "bobama@somewhere.com",
			after:            false,
			pgSize:           2,
			wantSize:         2,
			wantError:        nil,
			wantContainEmail: "alice@example.com",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

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
