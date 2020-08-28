package services_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/user"
	uuid "github.com/satori/go.uuid"
)

func TestUserCreate(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	tests := []struct {
		name      string
		username  string
		email     string
		password  string
		phone     string
		wantError []byte
	}{
		{"Valid user", "New User", "valid@user.com", "secret", "12345", nil},
		{"Bad uuid", "Bad uuid", "bad@uuid.com", "secret", "12345", []byte("repository ID not equal to new user ID")},
		{"DB error", "DB error", "db@error.com", "secret", "12345", []byte("DB error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.Create(tt.username, tt.email, tt.password, tt.phone)

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

func TestFindUserByID(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

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
			wantError: user.ErrNoRecord,
		},
		{
			desc:      "Invalid ID",
			userID:    "d1700797-42d460cda46f2448",
			wantUser:  nil,
			wantError: user.ErrInvalidInputSyntax,
		},
		{
			desc:      "empty ID",
			userID:    "",
			wantUser:  nil,
			wantError: user.ErrInvalidInputSyntax,
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

func TestFindUserByEmail(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

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
			wantError: user.ErrNoRecord,
		},
		{
			desc:      "empty Email",
			userEmail: "",
			wantUser:  nil,
			wantError: user.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			u, err := svc.FindByEmail(tC.userEmail)

			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if !reflect.DeepEqual(u, tC.wantUser) {
				t.Errorf("want \n%v\n; got \n%v\n", tC.wantUser, u)
			}

		})
	}
}

func TestUserAuthenticate(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	testCases := []struct {
		desc      string
		email     string
		pw        string
		wantUser  *user.User
		wantError error
	}{
		{
			desc:  "Valid User",
			email: "bob@example.com",
			pw:    "test1234",
			wantUser: &user.User{
				ID:                "68b1d5e2-39dd-4713-8631-a08100383a0f",
				Name:              "Bob",
				Email:             "bob@example.com",
				Phone:             "6544334535",
				HashedPw:          []byte("$2a$12$kHna5vstSSusP9VFC89tZ.317kInW7dZRL8snvnAej66wgQnyaQte"),
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				Active:            true,
				EmailWasValidated: true,
				Roles:             []string{user.RoleCommon, user.RoleAdmin},
			},
			wantError: nil,
		},
		{
			desc:  "Invalid Credential",
			email: "bob@example.com",
			pw:    "wrongpassword1234",
			wantUser: &user.User{
				ID:                "68b1d5e2-39dd-4713-8631-a08100383a0f",
				Name:              "Bob",
				Email:             "bob@example.com",
				Phone:             "6544334535",
				HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				Active:            true,
				EmailWasValidated: true,
				Roles:             []string{user.RoleCommon, user.RoleAdmin},
			},
			wantError: user.ErrInvalidCredentials,
		},
		{
			desc:      "Non-existing user",
			email:     "joe@nowhere.com",
			pw:        "test1234",
			wantUser:  nil,
			wantError: user.ErrNoRecord,
		},
		{
			desc:      "empty Email",
			email:     "",
			pw:        "test1234",
			wantUser:  nil,
			wantError: user.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			u, err := svc.Authenticate(tC.email, tC.pw)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				user, err := svc.FindByID(*u)
				if err != nil {
					t.Error("Could not find user")
				}
				if user.ID != tC.wantUser.ID {
					t.Errorf("want \n%v\n; got \n%v\n", tC.wantUser.ID, user.ID)
				}
				if user.Email != tC.wantUser.Email {
					t.Errorf("want \n%v\n; got \n%v\n", tC.wantUser.Email, user.Email)
				}
				if user.Name != tC.wantUser.Name {
					t.Errorf("want \n%v\n; got \n%v\n", tC.wantUser.Name, user.Name)
				}
				if user.Phone != tC.wantUser.Phone {
					t.Errorf("want \n%v\n; got \n%v\n", tC.wantUser.Phone, user.Phone)
				}
			}
		})
	}
}

func TestUpdateRoles(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	testCases := []struct {
		desc         string
		uid          string
		currentRoles []string
		newRoles     []string
		wantError    error
	}{
		{
			desc:         "Valid User",
			uid:          "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			currentRoles: []string{user.RoleCommon},
			newRoles:     []string{user.RoleAdmin},
			wantError:    nil,
		},
		{desc: "Invalid User",
			uid:          "85f45ff9-d31c-4ff7-d31c-5afb5a1f0fcd",
			currentRoles: []string{user.RoleCommon},
			newRoles:     []string{user.RoleAdmin},
			wantError:    user.ErrNoRecord,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.UpdateRoles(tC.uid, tC.newRoles)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				user, err := svc.FindByID(tC.uid)
				if err != nil {
					t.Error("Could not find user")
				}
				if reflect.DeepEqual(user.Roles, tC.currentRoles) {
					t.Errorf("want \n%v\n to be different from \n%v\n", tC.currentRoles, user.Roles)
				}
			}
		})
	}
}

func TestUpdatePw(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	testCases := []struct {
		desc      string
		email     string
		currentPw string
		newPw     string
		wantError error
	}{
		{
			desc:      "Valid User",
			email:     "bob@example.com",
			currentPw: "test1234",
			newPw:     "newPwtest1234",
			wantError: nil,
		},
		{
			desc:      "Empty new password",
			email:     "bob@example.com",
			currentPw: "newPwtest1234",
			newPw:     "",
			wantError: user.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// authenticate user
			uid, err := svc.Authenticate(tC.email, tC.currentPw)
			if err != nil {
				t.Error("Could not authenticate user")
			}
			if uid == nil {
				t.Error("returned uid = nil")
			} else {
				// update Pw
				err = svc.UpdatePw(*uid, tC.newPw)
				if err != tC.wantError {
					t.Errorf("want \n%v\n; got \n%v\n", tC.wantError, err)
				}
			}
			if tC.wantError == nil {
				// authenticate again
				uid, err = svc.Authenticate(tC.email, tC.newPw)
				if err != nil {
					t.Error("Could not authenticate with new password")
				}
			}
		})
	}
}

func TestUpdateUserStatus(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	testCases := []struct {
		desc      string
		uid       string
		newStatus bool
		wantError error
	}{
		{
			desc:      "Valid User",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newStatus: false,
			wantError: nil,
		},
		//
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.UpdateStatus(tC.uid, tC.newStatus)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				user, err := svc.FindByID(tC.uid)
				if err != nil {
					t.Error("Could not find user")
				}
				if (user != nil) && (user.Active != tC.newStatus) {
					t.Errorf("want \n%v\n; got \n%v\n", tC.newStatus, user.Active)
				}
			}
		})
	}
}

func TestUpdateEmailWasValidated(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	testCases := []struct {
		desc              string
		uid               string
		EmailWasValidated bool
		wantError         error
	}{
		{
			desc:              "Valid User",
			uid:               "68b1d5e2-39dd-4713-8631-a08100383a0f",
			EmailWasValidated: false,
			wantError:         nil,
		},
		//
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.UpdateEmailValidated(tC.uid, tC.EmailWasValidated)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				user, err := svc.FindByID(tC.uid)
				if err != nil {
					t.Error("Could not find user")
				}
				if (user != nil) && (user.EmailWasValidated != tC.EmailWasValidated) {
					t.Errorf("want \n%v\n; got \n%v\n", tC.EmailWasValidated, user.EmailWasValidated)
				}
			}
		})
	}
}

func TestUpdateUserName(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

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
			wantError:   user.ErrInvalidInputSyntax,
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

func TestUpdateEmail(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	testCases := []struct {
		desc      string
		uid       string
		newEmail  string
		wantError error
	}{
		{
			desc:      "Valid User",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newEmail:  "bob@newhost.com",
			wantError: nil,
		},
		{
			desc:      "Empty new user name",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newEmail:  "",
			wantError: user.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.UpdateEmail(tC.uid, tC.newEmail)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				user, err := svc.FindByID(tC.uid)
				if err != nil {
					t.Error("Could not find user")
				}
				if (user != nil) && (user.Email != tC.newEmail) {
					t.Errorf("want \n%v\n; got \n%v\n", tC.newEmail, user.Email)
				}
			}
		})
	}
}

func TestUpdatePhone(t *testing.T) {
	repo := mocks.NewUserRepo()
	svc := services.NewUserService(repo)

	testCases := []struct {
		desc      string
		uid       string
		newPhone  string
		wantError error
	}{
		{
			desc:      "Valid User",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newPhone:  "3453453452",
			wantError: nil,
		},
		{
			desc:      "Empty new user name",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newPhone:  "",
			wantError: user.ErrInvalidInputSyntax,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.UpdatePhone(tC.uid, tC.newPhone)
			if err != tC.wantError {
				t.Errorf("want %v; got %v", tC.wantError, err)
			}
			if tC.wantError == nil {
				user, err := svc.FindByID(tC.uid)
				if err != nil {
					t.Error("Could not find user")
				}
				if (user != nil) && (user.Phone != tC.newPhone) {
					t.Errorf("want \n%v\n; got \n%v\n", tC.newPhone, user.Phone)
				}
			}
		})
	}
}

func TestGetAll(t *testing.T) {
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
