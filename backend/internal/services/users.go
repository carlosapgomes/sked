package services

import (
	"encoding/base64"
	"errors"
	"time"

	"carlosapgomes.com/gobackend/internal/user"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService provides implementation of user domain interface
type userService struct {
	repo user.Repository
}

// NewUserService returns an user Service instance
func NewUserService(repo user.Repository) user.Service {
	return &userService{
		repo,
	}
}

// AdminRole constant
// const AdminRole = "admin"

// UserRole constant
// const UserRole = "user"

// CreateUser creates a new user and returns its uuid
func (s *userService) Create(name, email, password, phone string) (*string, error) {
	// // Create a bcrypt hash of the plain-text password.
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	uid := uuid.NewV4()
	dt := time.Now().UTC()
	roles := []string{}
	newUser := user.User{
		ID:                uid.String(),
		Name:              name,
		Email:             email,
		Phone:             phone,
		HashedPw:          hashedPw,
		CreatedAt:         dt,
		UpdatedAt:         dt,
		Active:            true,
		EmailWasValidated: false,
		Roles:             roles,
	}
	var id *string
	id, err = s.repo.Create(newUser)

	if err != nil {
		return nil, err
	}
	if (id != nil) && (*id != newUser.ID) {
		return nil, errors.New("New user creation: returned repository ID not equal to new user ID")
	}
	return id, err
}

// FindByID searches a user by its ID
func (s *userService) FindByID(uID string) (*user.User, error) {
	if uID == "" {
		return nil, user.ErrInvalidInputSyntax
	}
	_, err := uuid.FromString(uID)
	if err != nil {
		return nil, user.ErrInvalidInputSyntax
	}

	u, err := s.repo.FindByID(uID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail searches a user by its email address
func (s *userService) FindByEmail(email string) (*user.User, error) {
	if email == "" {
		return nil, user.ErrInvalidInputSyntax
	}

	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Authenticate autheticates a user and return its ID
func (s *userService) Authenticate(email, password string) (*string, error) {
	if (email == "") || (password == "") {
		return nil, user.ErrInvalidInputSyntax
	}
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(u.HashedPw, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, user.ErrInvalidCredentials
		}
		return nil, err
	}
	return &u.ID, nil
}

// UpdateRoles set a new user roles
func (s *userService) UpdateRoles(id string, roles []string) error {
	return s.repo.UpdateRoles(id, roles)
}

// UpdatePw updates a user's password
func (s *userService) UpdatePw(id, pw string) error {
	if pw == "" {
		return user.ErrInvalidInputSyntax
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	if err != nil {
		return err
	}
	return s.repo.UpdatePw(id, hashedPw)
}

// UpdateStatus updates a user's status
func (s *userService) UpdateStatus(id string, active bool) error {
	return s.repo.UpdateStatus(id, active)
}

// UpdateEmailValidated updates a user's EmailValidated
func (s *userService) UpdateEmailValidated(id string, emailWasValidated bool) error {
	return s.repo.UpdateEmailValidated(id, emailWasValidated)
}

// UpdateName, updates user name
func (s *userService) UpdateName(id string, name string) error {
	if name == "" {
		return user.ErrInvalidInputSyntax
	}
	return s.repo.UpdateName(id, name)
}

// UpdateEmail, updates user email, only admin should be able to do this.
// email account is the source of truth for system access
func (s *userService) UpdateEmail(id string, email string) error {
	// validate email string
	if email == "" {
		return user.ErrInvalidInputSyntax
	}
	return s.repo.UpdateEmail(id, email)
}

// UpdatePhone updates user email
func (s *userService) UpdatePhone(id string, phone string) error {
	if phone == "" {
		return user.ErrInvalidInputSyntax
	}
	return s.repo.UpdatePhone(id, phone)
}

// GetAll returns a paginated list of all users ordered by email
func (s *userService) GetAll(before string, after string, pgSize int) (*user.Cursor, error) {
	var usersResp user.Cursor
	var err error
	var uList *[]user.User
	switch {
	case (before != "" && after != ""):
		// if both (before & after) are present, returns error
		return nil, user.ErrInvalidInputSyntax
	case (before == "" && after == ""):
		// if they are empty/or absent
		// get default list and page size
		uList, usersResp.HasBefore, err = s.repo.
			GetAll("", false, pgSize)
		if err != nil {
			return nil, err
		}
		if uList != nil {
			for _, u := range *uList {
				usersResp.Users = append(usersResp.Users, u)
			}
		}
		if len(usersResp.Users) > 0 {
			usersResp.Before = base64.StdEncoding.
				EncodeToString([]byte(usersResp.Users[len(usersResp.Users)-1].Email))
		} else {
			usersResp.Before = ""
		}
		usersResp.After = ""
		usersResp.HasAfter = false
		// and return values
	case (before != ""):
		// if before is present,
		// get a before list
		c, err := base64.StdEncoding.DecodeString(before)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		uList, usersResp.HasBefore, err = s.repo.GetAll(cursor, false, pgSize)
		if err != nil {
			return nil, err
		}
		if uList != nil {
			for _, u := range *uList {
				usersResp.Users = append(usersResp.Users, u)
			}
		}
		if len(usersResp.Users) > 0 {
			befCursor := base64.StdEncoding.EncodeToString([]byte(usersResp.Users[len(usersResp.Users)-1].Email))
			usersResp.Before = befCursor
		} else {
			usersResp.Before = ""
		}
		// test for 'after data' from the requested cursor
		// fill the response fields
		_, usersResp.HasAfter, err = s.repo.GetAll(cursor, true, pgSize)
		if usersResp.HasAfter {
			usersResp.After = base64.StdEncoding.EncodeToString([]byte(before))
		} else {
			usersResp.After = ""
		}
		// and return it
	case (after != ""):
		// if after is present,
		// get an after list
		c, err := base64.StdEncoding.DecodeString(after)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		uList, usersResp.HasAfter, err = s.repo.
			GetAll(cursor, true, pgSize)
		// and return it
		if uList != nil {
			for _, u := range *uList {
				usersResp.Users = append(usersResp.Users, u)
			}
		}
		if len(usersResp.Users) > 0 {
			usersResp.After = base64.StdEncoding.EncodeToString([]byte(usersResp.Users[0].Email))
		}
		// test for 'before data' from the requested cursor
		// fill the response fields
		_, usersResp.HasBefore, err = s.repo.
			GetAll(cursor, false, pgSize)
		if usersResp.HasBefore {
			usersResp.Before = base64.StdEncoding.EncodeToString([]byte(after))
		}
	}
	return &usersResp, nil
}

// FindByName returns a list of users whose names looks like 'name'
func (s *userService) FindByName(name string, before string, hasBef bool, after string, hasAft bool, pgSize int) (*[]user.User, error) {
	return nil, nil
}
