package mocks

import (
	"strings"
	"time"

	"carlosapgomes.com/sked/internal/user"
	"github.com/pkg/errors"
)

// mocks user services and repository

// UserMockRepo is a mocked user repository
type UserMockRepo struct {
	uDb []user.User
}

//NewUserRepo returns a mocked repository
func NewUserRepo() *UserMockRepo {
	var db []user.User
	validUser := &user.User{
		ID:                "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		Name:              "Valid User",
		Email:             "valid@user.com",
		Phone:             "6544332135",
		HashedPw:          []byte("validPw"),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Active:            true,
		EmailWasValidated: true,
		Roles:             []string{user.RoleClerk},
	}
	validUser2 := &user.User{
		ID:                "68b1d5e2-39dd-4713-8631-a08100383a0f",
		Name:              "Bob",
		Email:             "bob@example.com",
		Phone:             "6544334535",
		HashedPw:          []byte("$2a$12$kHna5vstSSusP9VFC89tZ.317kInW7dZRL8snvnAej66wgQnyaQte"),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Active:            true,
		EmailWasValidated: true,
		Roles:             []string{user.RoleClerk, user.RoleAdmin},
	}
	validUser3 := &user.User{
		ID:                "dcce1beb-aee6-4a4d-b724-94d470817323",
		Name:              "Alice Jones",
		Email:             "alice@example.com",
		Phone:             "6544332135",
		HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
		CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Active:            true,
		EmailWasValidated: true,
		Roles:             []string{user.RoleClerk},
	}
	validUser4 := &user.User{
		ID:                "ecadbb28-14e6-4560-8574-809c6c54b9cb",
		Name:              "Barack Obama",
		Email:             "bobama@somewhere.com",
		Phone:             "6544332135",
		HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
		CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Active:            false,
		EmailWasValidated: true,
		Roles:             []string{user.RoleClerk},
	}
	validUser5 := &user.User{
		ID:                "ca16fc9d-df7b-4594-97e3-264432145b01",
		Name:              "SpongeBob Squarepants",
		Email:             "spongebob@somewhere.com",
		Phone:             "65949340",
		HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
		CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Active:            false,
		EmailWasValidated: true,
		Roles:             []string{user.RoleClerk},
	}
	validUser6 := &user.User{
		ID:                "27f9802b-acb3-4852-bf97-c4ed4c3b3658",
		Name:              "Tim Berners-Lee",
		Email:             "tblee@somewhere.com",
		Phone:             "0323949324",
		HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
		CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Active:            false,
		EmailWasValidated: true,
		Roles:             []string{user.RoleClerk},
	}
	validUser7 := &user.User{
		ID:                "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Name:              "Dr. House",
		Email:             "house@doctor.com",
		Phone:             "6544332135",
		HashedPw:          []byte("validPw"),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Active:            true,
		EmailWasValidated: true,
		Roles:             []string{user.RoleDoctor},
	}
	validUser8 := &user.User{
		ID:                "a520df95-02fa-4d86-8eef-58385c354b29",
		Name:              "Shaun Murphy",
		Email:             "shaun@thegooddoctor.com",
		Phone:             "64532332135",
		HashedPw:          []byte("validPw"),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Active:            true,
		EmailWasValidated: true,
		Roles:             []string{user.RoleDoctor},
	}
	db = append(db, *validUser)
	db = append(db, *validUser2)
	db = append(db, *validUser3)
	db = append(db, *validUser4)
	db = append(db, *validUser5)
	db = append(db, *validUser6)
	db = append(db, *validUser7)
	db = append(db, *validUser8)
	return &UserMockRepo{
		db,
	}
}

// Create mocks user creation
func (r *UserMockRepo) Create(user user.User) (*string, error) {
	badID := "12342342"
	switch user.Name {
	case "New User":
		return &user.ID, nil
	case "Bad uuid":
		return &badID, nil
	case "DB error":
		return nil, errors.New("DB error")
	default:
		return nil, nil
	}
}

// UpdateStatus mocks updating user's status
func (r *UserMockRepo) UpdateStatus(id string, active bool) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].Active = active
			return nil
		}
	}
	return user.ErrNoRecord
}

// UpdateEmailValidated mocks updating user's status
func (r *UserMockRepo) UpdateEmailValidated(id string, emailWasValidated bool) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].EmailWasValidated = emailWasValidated
			return nil
		}
	}
	return user.ErrNoRecord
}

// UpdateRoles mocks updating user's roles
func (r *UserMockRepo) UpdateRoles(id string, roles []string) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].Roles = roles
			return nil
		}
	}
	return user.ErrNoRecord
}

// UpdateEmail mocks updating user's Email
func (r *UserMockRepo) UpdateEmail(id string, email string) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].Email = email
			return nil
		}
	}
	return user.ErrNoRecord
}

// UpdateName mocks updating user's Name
func (r *UserMockRepo) UpdateName(id string, name string) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].Name = name
			return nil
		}
	}
	return user.ErrNoRecord
}

// UpdatePhone mocks updating user's Phone
func (r *UserMockRepo) UpdatePhone(id string, phone string) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].Phone = phone
			return nil
		}
	}
	return user.ErrNoRecord
}

// UpdatePw mocks updating user's Pw
func (r *UserMockRepo) UpdatePw(id string, pw []byte) error {
	for i, u := range r.uDb {
		if u.ID == id {
			r.uDb[i].HashedPw = pw
			return nil
		}
	}
	return user.ErrNoRecord
}

// FindByID mocks finding a user by its ID
func (r *UserMockRepo) FindByID(id string) (*user.User, error) {
	for _, u := range r.uDb {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, user.ErrNoRecord
}

// FindByEmail mocks finding a user by its ID
func (r *UserMockRepo) FindByEmail(email string) (*user.User, error) {
	for _, u := range r.uDb {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, user.ErrNoRecord
}

// GetAll returns a list of users ordered by email
func (r *UserMockRepo) GetAll(cursor string, after bool, pgSize int) (*[]user.User, bool, error) {
	if cursor == "" {
		return &r.uDb, false, nil
	}
	pos := r.findPos(r.uDb, cursor)
	if pos == -1 {
		return nil, false, user.ErrNoRecord
	}

	var res []user.User
	var hasMore bool
	hasMore = false
	if after {
		start := pos + 1
		for i := start; i < (start + pgSize); i++ {
			res = append(res, r.uDb[i])
		}
		if (len(r.uDb) - pos) > pgSize {
			hasMore = true
		}
	} else {
		start := pos - pgSize
		if start < 0 {
			start = 0
		}
		for i := start; i <= (pos - 1); i++ {
			res = append(res, r.uDb[i])
		}
		if pos > pgSize {
			hasMore = true
		}
	}
	return &res, hasMore, nil
}

func (r UserMockRepo) findPos(patients []user.User, id string) int {
	for i, el := range patients {
		if el.ID == id {
			return i
		}
	}
	return -1
}

// FindByName returns a list of users whose names looks like 'name'
func (r *UserMockRepo) FindByName(name string) (*[]user.User, error) {
	return nil, nil
}

// UserMockSvc mocks user services
type UserMockSvc struct {
}

// NewUserSvc returns a mocked user service
func NewUserSvc() *UserMockSvc {
	return &UserMockSvc{}
}

// Create mocks new user creation service
func (s UserMockSvc) Create(name, email, password, phone string) (*string, error) {
	switch email {
	case "dupe@example.com":
		return nil, errors.Wrap(user.ErrDuplicateField, "email already in use")
	default:
		uid := "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd"
		return &uid, nil
	}
}

// FindByID mocks finding a user by its ID
func (s UserMockSvc) FindByID(id string) (*user.User, error) {
	if id == "68b1d5e2-39dd-4713-8631-a08100383a0f" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &user.User{
			ID:                "68b1d5e2-39dd-4713-8631-a08100383a0f",
			Name:              "Bob",
			Email:             "bob@example.com",
			Phone:             "6544334535",
			HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"), // pw: test1234
			CreatedAt:         dt,
			UpdatedAt:         dt,
			Active:            true,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk, user.RoleAdmin},
		}, nil
	}
	if id == "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &user.User{
			ID:                "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			Name:              "Valid User",
			Email:             "valid@user.com",
			Phone:             "6544332135",
			HashedPw:          []byte("validPw"),
			CreatedAt:         dt,
			UpdatedAt:         dt,
			Active:            true,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk},
		}, nil
	}
	if id == "dcce1beb-aee6-4a4d-b724-94d470817323" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &user.User{
			ID:                "dcce1beb-aee6-4a4d-b724-94d470817323",
			Name:              "Alice Jones",
			Email:             "alice@example.com",
			Phone:             "6544332135",
			HashedPw:          []byte("validPw"),
			CreatedAt:         dt,
			UpdatedAt:         dt,
			Active:            true,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk},
		}, nil
	}
	if id == "ecadbb28-14e6-4560-8574-809c6c54b9cb" {
		return &user.User{
			ID:                "ecadbb28-14e6-4560-8574-809c6c54b9cb",
			Name:              "Barack Obama",
			Email:             "bobama@somewhere.com",
			Phone:             "6544332135",
			HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
			CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			Active:            false,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk},
		}, nil
	}
	return nil, errors.New("not found")
}

// FindByEmail mocks finding a user by its ID
func (s UserMockSvc) FindByEmail(email string) (*user.User, error) {
	if email == "bob@example.com" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &user.User{
			ID:                "68b1d5e2-39dd-4713-8631-a08100383a0f",
			Name:              "Bob",
			Email:             "bob@example.com",
			Phone:             "6544334535",
			HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"), // pw: test1234
			CreatedAt:         dt,
			UpdatedAt:         dt,
			Active:            true,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk, user.RoleAdmin},
		}, nil
	}
	if email == "valid@user.com" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &user.User{
			ID:                "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
			Name:              "Valid User",
			Email:             "valid@user.com",
			Phone:             "6544332135",
			HashedPw:          []byte("validPw"),
			CreatedAt:         dt,
			UpdatedAt:         dt,
			Active:            true,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk},
		}, nil
	}
	if email == "alice@example.com" {
		dt := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		return &user.User{
			ID:                "dcce1beb-aee6-4a4d-b724-94d470817323",
			Name:              "Alice Jones",
			Email:             "alice@example.com",
			Phone:             "6544332135",
			HashedPw:          []byte("validPw"),
			CreatedAt:         dt,
			UpdatedAt:         dt,
			Active:            true,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk},
		}, nil
	}
	if email == "bobama@somewhere.com" {
		return &user.User{
			ID:                "ecadbb28-14e6-4560-8574-809c6c54b9cb",
			Name:              "Barack Obama",
			Email:             "bobama@somewhere.com",
			Phone:             "6544332135",
			HashedPw:          []byte("$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je"),
			CreatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			Active:            false,
			EmailWasValidated: true,
			Roles:             []string{user.RoleClerk},
		}, nil
	}
	return nil, user.ErrNoRecord
}

// TODO: should we have a RemoveRole?

// UpdateRoles updates user roles
func (s UserMockSvc) UpdateRoles(id string, roles []string) error {
	return nil
}

// UpdatePw mocks user's password updating
func (s UserMockSvc) UpdatePw(id, password string) error {
	return nil
}

// UpdateStatus mocks updating user's status
func (s UserMockSvc) UpdateStatus(id string, active bool) error {
	return nil
}

// UpdateEmailValidated mocks updating user's status
func (s UserMockSvc) UpdateEmailValidated(id string, active bool) error {
	return nil
}

// UpdateName updates user name
func (s UserMockSvc) UpdateName(id string, name string) error {
	return nil
}

// UpdateEmail updates user email
func (s UserMockSvc) UpdateEmail(id string, email string) error {
	return nil
}

// UpdatePhone updates user email
func (s UserMockSvc) UpdatePhone(id string, phone string) error {
	return nil
}

// Authenticate verify user's credentials and returns user ID
func (s UserMockSvc) Authenticate(email, password string) (*string, error) {
	if (email == "bob@example.com") && (password == "test1234") {
		id := "68b1d5e2-39dd-4713-8631-a08100383a0f"
		return &id, nil
	}
	return nil, user.ErrInvalidCredentials
}

// GetAll return a lista of users ordered by email
func (s UserMockSvc) GetAll(before string, after string, pgSize int) (*user.Cursor, error) {
	return nil, nil
}

// FindByName return a list of users whose names looks like 'name'
func (s UserMockSvc) FindByName(name string) (*[]user.User, error) {
	repo := NewUserRepo()
	var users []user.User
	for _, u := range repo.uDb {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(name)) {
			users = append(users, u)
		}
	}
	if len(users) == 0 {
		return nil, user.ErrNoRecord
	}
	return &users, nil
}
