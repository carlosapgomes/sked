package storage_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	"carlosapgomes.com/gobackend/internal/storage"
	"carlosapgomes.com/gobackend/internal/user"
	"golang.org/x/crypto/bcrypt"
)

func TestCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name      string
		newUser   *user.User
		wantError error
	}{
		{
			name: "Valid User",
			newUser: &user.User{
				ID:        "5b28f739-e372-4622-8390-9992f3c7b0e9",
				Name:      "Muhamed Ali",
				Email:     "muhamedali@nowhere.com",
				HashedPw:  []byte("$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG"),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Active:    true,
				Roles:     []string{},
			},
			wantError: nil,
		},
		{
			name: "Bad userID",
			newUser: &user.User{
				ID:        "5b28f739-e372-4622-9992f3c7b0e9",
				Name:      "Muhamed Ali",
				Email:     "muhamedali@another.com",
				HashedPw:  []byte("$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG"),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Active:    true,
				Roles:     []string{},
			},
			wantError: user.ErrInvalidInputSyntax,
		},
		{
			name: "Duplicate ID",
			newUser: &user.User{
				ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
				Name:      "Alice Jones",
				Email:     "alice@example.com",
				HashedPw:  []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Active:    true,
				Roles:     []string{},
			},
			wantError: user.ErrDuplicateField,
		},
		{
			name: "Duplicate email",
			newUser: &user.User{
				ID:        "75499ef5-bde6-4f39-81b8-daf181942741",
				Name:      "Alice Jones",
				Email:     "alice@example.com",
				HashedPw:  []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Active:    true,
				Roles:     []string{},
			},
			wantError: user.ErrDuplicateField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)

			uid, err := repo.Create(*tt.newUser)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if (uid == nil) && (tt.wantError == nil) {
				t.Errorf("want %s; got uid = nil", tt.newUser.ID)
			}
			if uid != nil {
				user, _ := repo.FindByID(*uid)
				if (tt.newUser.Name != user.Name) || (tt.newUser.Email != user.Email) {
					t.Errorf("want \n%v\n; got \n%v\n", tt.newUser.Name, user.Name)
				}
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    string
		wantUser  *user.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: "dcce1beb-aee6-4a4d-b724-94d470817323",
			wantUser: &user.User{
				ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
				Name:      "Alice Jones",
				Email:     "alice@example.com",
				Phone:     "6544334535",
				HashedPw:  []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
				CreatedAt: time.Date(2019, 06, 23, 17, 25, 00, 0, time.UTC),
				UpdatedAt: time.Date(2019, 06, 23, 17, 25, 00, 0, time.UTC),
				Active:    true,
				Roles:     []string{user.RoleCommon},
			},
			wantError: nil,
		},
		{
			name:      "Non-existing ID",
			userID:    "d1700797-42d4-4fe4-8fc2-60cda46f2448",
			wantUser:  nil,
			wantError: user.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)

			user, err := repo.FindByID(tt.userID)

			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want \n%v\n; got \n%v\n", tt.wantUser, user)
			}
		})
	}
}

func TestFindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name      string
		userEmail string
		wantUser  *user.User
		wantError error
	}{
		{
			name:      "Valid Email",
			userEmail: "alice@example.com",
			wantUser: &user.User{
				ID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
				Name:      "Alice Jones",
				Email:     "alice@example.com",
				Phone:     "6544334535",
				HashedPw:  []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
				CreatedAt: time.Date(2019, 06, 23, 17, 25, 00, 0, time.UTC),
				UpdatedAt: time.Date(2019, 06, 23, 17, 25, 00, 0, time.UTC),
				Active:    true,
				Roles:     []string{user.RoleCommon},
			},
			wantError: nil,
		},
		{
			name:      "Non-existing email",
			userEmail: "nobody@test.com",
			wantUser:  nil,
			wantError: user.ErrNoRecord,
		},
		{
			name:      "Empty email",
			userEmail: "",
			wantUser:  nil,
			wantError: user.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)

			user, err := repo.FindByEmail(tt.userEmail)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want \n%v\n; got \n%v\n", tt.wantUser, user)
			}
		})
	}
}

func TestUpdatePw(t *testing.T) {
	testCases := []struct {
		desc      string
		uid       string
		newPw     string
		wantError error
	}{
		{
			desc:      "Valid uid and new pw",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newPw:     "newTestPw1234",
			wantError: nil,
		},
		{
			desc:      "Invalid uid",
			uid:       "68b1d5e2-8631-a08100383a0f",
			newPw:     "newTestPw1234",
			wantError: errors.New("Any error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)
			hashedPw, _ := bcrypt.GenerateFromPassword([]byte(tC.newPw), 12)
			err := repo.UpdatePw(tC.uid, hashedPw)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			u, err := repo.FindByID(tC.uid)
			if u != nil && bytes.Compare(u.HashedPw, hashedPw) != 0 {
				t.Errorf("Want %v; got %v\n", hashedPw, u.HashedPw)
			}

		})
	}
}

func TestUpdateStatus(t *testing.T) {
	testCases := []struct {
		desc      string
		uid       string
		newActive bool
		wantError error
	}{
		{
			desc:      "Valid uid",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newActive: false,
			wantError: nil,
		},
		{
			desc:      "Invalid uid",
			uid:       "68b1d5e2-8631-a08100383a0f",
			newActive: false,
			wantError: errors.New("Any error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)
			err := repo.UpdateStatus(tC.uid, tC.newActive)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			u, err := repo.FindByID(tC.uid)
			if u != nil && u.Active != tC.newActive {
				t.Errorf("Want %v; got %v\n", tC.newActive, u.Active)
			}

		})
	}
}

func TestUpdateEmailValidated(t *testing.T) {
	testCases := []struct {
		desc                 string
		uid                  string
		newEmailWasValidated bool
		wantError            error
	}{
		{
			desc:                 "Valid uid",
			uid:                  "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newEmailWasValidated: false,
			wantError:            nil,
		},
		{
			desc:                 "Invalid uid",
			uid:                  "68b1d5e2-8631-a08100383a0f",
			newEmailWasValidated: false,
			wantError:            errors.New("Any error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)
			err := repo.UpdateEmailValidated(tC.uid, tC.newEmailWasValidated)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			u, err := repo.FindByID(tC.uid)
			if u != nil && u.EmailWasValidated != tC.newEmailWasValidated {
				t.Errorf("Want %v; got %v\n", tC.newEmailWasValidated, u.EmailWasValidated)
			}

		})
	}
}

func TestUpdateName(t *testing.T) {
	testCases := []struct {
		desc      string
		uid       string
		newName   string
		wantError error
	}{
		{
			desc:      "Valid uid",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newName:   "Johnny Smith",
			wantError: nil,
		},
		{
			desc:      "Invalid uid",
			uid:       "68b1d5e2-8631-a08100383a0f",
			newName:   "Johnny Smith",
			wantError: errors.New("Any error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)
			err := repo.UpdateName(tC.uid, tC.newName)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			u, err := repo.FindByID(tC.uid)
			if u != nil && u.Name != tC.newName {
				t.Errorf("Want %v; got %v\n", tC.newName, u.Name)
			}

		})
	}
}

func TestUpdatePhone(t *testing.T) {
	testCases := []struct {
		desc      string
		uid       string
		newPhone  string
		wantError error
	}{
		{
			desc:      "Valid uid",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newPhone:  "214377669988",
			wantError: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)
			u, _ := repo.FindByID(tC.uid)
			t.Log(u.Phone)
			t.Log(u.Name)
			t.Log(u.ID)
			t.Log(u.Roles)
			t.Log(u.Email)
			t.Log(u.EmailWasValidated)
			err := repo.UpdatePhone(tC.uid, tC.newPhone)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			u, err = repo.FindByID(tC.uid)
			if u != nil && u.Phone != tC.newPhone {
				t.Errorf("Want %v; got %v\n %v\n %v\n", tC.newPhone, u.Phone, u.Name, u.ID)
			}

		})
	}
}

func TestUpdateEmail(t *testing.T) {
	testCases := []struct {
		desc      string
		uid       string
		newEmail  string
		wantError error
	}{
		{
			desc:      "Valid uid",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newEmail:  "bobs@newhost.com",
			wantError: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)

			err := repo.UpdateEmail(tC.uid, tC.newEmail)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			u, err := repo.FindByID(tC.uid)
			if u != nil && u.Email != tC.newEmail {
				t.Errorf("Want %v; got %v\n", tC.newEmail, u.Email)
			}

		})
	}
}

func TestUpdateRoles(t *testing.T) {
	testCases := []struct {
		desc      string
		uid       string
		newRoles  []string
		wantError error
	}{
		{
			desc:      "Valid uid",
			uid:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
			newRoles:  []string{user.RoleAdmin},
			wantError: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgUserRepository(db)

			err := repo.UpdateRoles(tC.uid, tC.newRoles)
			if tC.wantError != nil && err == nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			if tC.wantError == nil && err != nil {
				t.Errorf("Want %v; got %v\n", tC.wantError, err)
			}
			u, err := repo.FindByID(tC.uid)
			if (u != nil) && !(reflect.DeepEqual(u.Roles, tC.newRoles)) {
				t.Errorf("Want %v; got %v\n", tC.newRoles, u.Roles)
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
