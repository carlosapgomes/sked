package services

import (
	"encoding/base64"
	"errors"
	"time"

	"carlosapgomes.com/sked/internal/user"
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

// GetAllDoctors returns a list of doctors ordered by name
func (s *userService) GetAllDoctors() (*[]user.User, error) {
	docs, err := s.repo.GetAllDoctors()
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// GetAll returns a paginated list of all users ordered by email
func (s *userService) GetAll(previous string, next string, pgSize int) (*user.Page, error) {
	var page user.Page
	var err error
	var list *[]user.User
	if pgSize <= 0 {
		return nil, user.ErrInvalidInputSyntax
	}
	switch {
	case (previous != "" && next != ""):
		// if both (previous & next) are present, return error
		return nil, user.ErrInvalidInputSyntax
	case (previous == "" && next == ""):
		// if both are empty, get the first "pgSize" elements of the list
		list, page.HasNextPage, err = s.repo.
			GetAll("", true, pgSize)
		if err != nil {
			return nil, err
		}
		if list != nil {
			for _, u := range *list {
				page.Users = append(page.Users, u)
			}
		}
		if len(page.Users) > 0 {
			page.StartCursor = base64.StdEncoding.
				EncodeToString([]byte(page.Users[0].ID))
			page.EndCursor = base64.StdEncoding.
				EncodeToString([]byte(page.
					Users[len(page.Users)-1].ID))
		} else {
			page.StartCursor = ""
			page.EndCursor = ""
		}
		page.HasPreviousPage = false
	case (previous != ""):
		// if previous is present, get a previous list
		c, err := base64.StdEncoding.DecodeString(previous)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		list, page.HasPreviousPage, err = s.repo.GetAll(cursor, false, pgSize)
		if err != nil {
			return nil, err
		}
		// test if list is not empty
		if list != nil && len(*list) > 0 {
			for _, item := range *list {
				page.Users = append(page.Users, item)
			}
		}
		if len(page.Users) > 0 {
			//fmt.Printf("StartCursor: %v\n", page.Users[0].ID)
			page.StartCursor = base64.StdEncoding.
				EncodeToString([]byte(page.Users[0].ID))
			page.EndCursor = base64.StdEncoding.
				EncodeToString([]byte(page.
					Users[len(page.Users)-1].ID))
			page.HasNextPage = true
		} else {
			page.StartCursor = ""
			page.EndCursor = ""
			page.HasNextPage = false
		}
	case (next != ""):
		// if next is present,
		// get an next list
		c, err := base64.StdEncoding.DecodeString(next)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		list, page.HasNextPage, err = s.repo.
			GetAll(cursor, true, pgSize)
		if err != nil {
			return nil, err
		}
		// test if list is not empty
		if list != nil && len(*list) > 0 {
			for _, item := range *list {
				page.Users = append(page.Users, item)
			}
		}
		if len(page.Users) > 0 {
			page.StartCursor = base64.StdEncoding.
				EncodeToString([]byte(page.Users[0].ID))
			page.EndCursor = base64.StdEncoding.
				EncodeToString([]byte(page.
					Users[len(page.Users)-1].ID))
			page.HasPreviousPage = true
		} else {
			page.StartCursor = ""
			page.EndCursor = ""
			page.HasPreviousPage = false
		}
	}
	return &page, nil
}

// FindByName returns a list of users whose names looks like 'name'
func (s *userService) FindByName(name string) (*[]user.User, error) {
	res, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	return res, nil
}
