package services

import (
	"time"

	"carlosapgomes.com/sked/internal/token"
	uuid "github.com/satori/go.uuid"
)

type tokenService struct {
	repo token.Repository
}

// NewTokenService returns a token service instance
func NewTokenService(repo token.Repository) token.Service {
	return &tokenService{
		repo,
	}
}

// New create a new token
func (s *tokenService) New(uid string, kind string) (*string, error) {
	dt := time.Now().UTC()
	// token length = 40 random characters
	id := newLen(40)
	var expiresIn time.Duration
	switch kind {
	case token.PwReset:
		expiresIn = token.ExpPwReset
	case token.ValidateEmail:
		expiresIn = token.ExpValidateEmail
	default:
		return nil, token.ErrInvalidInputSyntax
	}
	// Check id format
	if !s.isValidUUID(uid) {
		return nil, token.ErrInvalidInputSyntax
	}

	t := &token.Token{
		ID:        id,
		UID:       uid,
		CreatedAt: dt,
		ExpiresAt: dt.Add(expiresIn),
		Kind:      kind,
	}
	tID, err := s.repo.Save(t)
	if err != nil {
		return nil, err
	}
	if *tID != t.ID {
		return nil, token.ErrDb
	}
	return &id, nil
}

// FindByID find a token by its ID
func (s *tokenService) FindByID(id string) (*token.Token, error) {
	token, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// FindByUID find all tokens belonging to a user
func (s *tokenService) FindAllByUID(uid string) (*[]token.Token, error) {
	if !s.isValidUUID(uid) {
		return nil, token.ErrInvalidInputSyntax
	}
	tokens, err := s.repo.FindAllByUID(uid)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Delete erases a token
func (s *tokenService) Delete(id string) error {
	if id == "" {
		return token.ErrInvalidInputSyntax
	}
	return s.repo.Delete(id)
}

func (s *tokenService) isValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return (err == nil)
}
