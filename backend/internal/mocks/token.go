package mocks

import (
	"crypto/rand"
	"time"

	"carlosapgomes.com/sked/internal/token"
)

// Generates random string for a token ID
// From https://github.com/dchest/uniuri

// stdChars is a set of standard characters allowed in uniuri string.
var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// NewLen returns a new random string of the provided length, consisting of
// standard characters.
func newLen(length int) string {
	return newLenChars(length, stdChars)
}

// newLenChars returns a new random string of the provided length, consisting
// of the provided byte slice of allowed characters (maximum 256).
func newLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("uniuri: wrong charset length for newLenChars")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

// TokenMockSvc mocks a token service
type TokenMockSvc struct{}

// NewTokenMockSvc creates a new token service
func NewTokenMockSvc() *TokenMockSvc {
	return &TokenMockSvc{}
}

// New mocks creating a new token
func (s *TokenMockSvc) New(uid string, kind string) (*string, error) {
	tokenID := newLen(40)
	return &tokenID, nil
}

// FindByID mocks finding a token by its ID
func (s *TokenMockSvc) FindByID(id string) (*token.Token, error) {
	switch id {
	case "":
		return nil, token.ErrInvalidInputSyntax
	case "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2":
		dt := time.Now().UTC()
		return &token.Token{
			ID:        "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2",
			UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
			CreatedAt: dt,
			ExpiresAt: dt.Add(time.Duration(30) * time.Minute),
			Kind:      "PwReset",
		}, nil
	case "U5UC7FXsSq6I9632BndMkSO4i1Z83v9KGd2FDCDN":
		// expired token
		dt := time.Now().UTC()
		creatT := dt.Add(time.Duration(-60) * time.Minute)
		return &token.Token{
			ID:        "U5UC7FXsSq6I9632BndMkSO4i1Z83v9KGd2FDCDN",
			UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
			CreatedAt: creatT,
			ExpiresAt: creatT.Add(time.Duration(10) * time.Minute),
			Kind:      "PwReset",
		}, nil
	case "Sqi1Z86I9632BndMkS632BNO4i1Z83v9KGd27FXs":
		dt := time.Now().UTC()
		return &token.Token{
			ID:        "Sqi1Z86I9632BndMkS632BNO4i1Z83v9KGd27FXs",
			UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
			CreatedAt: dt,
			ExpiresAt: dt.Add(time.Duration(30) * time.Minute),
			Kind:      "ValidateEmail",
		}, nil
	case "6Lg3qdCsAgNDFgm7HNO2GUwhagEu2v7syMVLSpZ5":
		// expired token
		dt := time.Now().UTC()
		creatT := dt.Add(time.Duration(-60) * time.Minute)
		return &token.Token{
			ID:        "6Lg3qdCsAgNDFgm7HNO2GUwhagEu2v7syMVLSpZ5",
			UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
			CreatedAt: creatT,
			ExpiresAt: creatT.Add(time.Duration(10) * time.Minute),
			Kind:      "ValidateEmail",
		}, nil
	default:
		return nil, token.ErrNoRecord
	}
}

// FindAllByUID mocks finding all tokens belonging to a user
func (s *TokenMockSvc) FindAllByUID(uid string) (*[]token.Token, error) {
	if uid == "68b1d5e2-39dd-4713-8631-a08100383a0f" {
		dt := time.Now().UTC()
		return &[]token.Token{
			{
				ID:        "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2",
				UID:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
				CreatedAt: dt,
				ExpiresAt: dt.Add(time.Duration(30) * time.Minute),
				Kind:      "PwReset",
			},
			{
				ID:        "iJiY6njKPdfPxZksvv27vM7C5cM8myaGQ4rnAzwm",
				UID:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
				CreatedAt: dt,
				ExpiresAt: dt.Add(time.Duration(30) * time.Minute),
				Kind:      "PwReset",
			},
		}, nil
	}
	return nil, token.ErrNoRecord
}

// Delete mocks deleting a token
func (s *TokenMockSvc) Delete(id string) error {
	return nil
}

// TokenMockRepo mocks a token repository
type TokenMockRepo struct{}

// NewTokenRepo creates a new token repository
func NewTokenRepo() *TokenMockRepo {
	return &TokenMockRepo{}
}

// FindByID mocks finding a token by its ID
func (r *TokenMockRepo) FindByID(id string) (*token.Token, error) {
	switch id {
	case "":
		return nil, token.ErrInvalidInputSyntax
	case "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2":
		dt := time.Now().UTC()
		return &token.Token{
			ID:        "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2",
			UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
			CreatedAt: dt,
			ExpiresAt: dt.Add(time.Duration(30) * time.Minute),
			Kind:      "PwReset",
		}, nil
	default:
		return nil, token.ErrNoRecord
	}
}

// FindAllByUID mocks finding all token belonging to a user
func (r *TokenMockRepo) FindAllByUID(uid string) (*[]token.Token, error) {
	if uid == "dcce1beb-aee6-4a4d-b724-94d470817323" {
		dt := time.Now().UTC()
		return &[]token.Token{
			{
				ID:        "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: dt,
				ExpiresAt: dt.Add(time.Duration(30) * time.Minute),
				Kind:      "PwReset",
			},
			{
				ID:        "iJiY6njKPdfPxZksvv27vM7C5cM8myaGQ4rnAzwm",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: dt,
				ExpiresAt: dt.Add(time.Duration(30) * time.Minute),
				Kind:      "PwReset",
			},
		}, nil
	}
	return nil, token.ErrNoRecord
}

// Save mocks persisting a token
func (r *TokenMockRepo) Save(t *token.Token) (*string, error) {
	switch t.UID {
	case "44bd844b-a0b9-4d48-a371-cc2a44a7a4e3":
		return nil, token.ErrDb
	case "fbf344b4-2b23-4de6-9541-9678b904c020":
		// return wrong ID
		id := "632BndMkSFDCDNO4i1Z83v9KGd2"
		return &id, nil
	default:
		return &t.ID, nil
	}
}

// Delete mocks a token deletion
func (r *TokenMockRepo) Delete(id string) error {
	if id == "iJiY6njKPdfPxZksvv27vM7C5cM8myaGQ4rnAzwm" {
		return nil
	}
	return token.ErrNoRecord
}
