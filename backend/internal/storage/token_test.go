package storage_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/storage"
	"carlosapgomes.com/sked/internal/token"
)

// TestFindByID tests finding a token by its ID
func TestFindTokenByID(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name      string
		tokenID   string
		wantToken *token.Token
		wantError error
	}{
		{
			name:    "valid ID",
			tokenID: "7telsDIFzlzIlZ1fiH8pDtoFiJMoBUi69j9525jt",
			wantToken: &token.Token{
				ID:        "7telsDIFzlzIlZ1fiH8pDtoFiJMoBUi69j9525jt",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Date(2019, 06, 23, 18, 25, 0, 0, time.UTC),
				ExpiresAt: time.Date(2019, 06, 23, 18, 55, 0, 0, time.UTC),
				Kind:      "PwReset",
			},
			wantError: nil,
		},
		{
			name:      "invalid ID",
			tokenID:   "",
			wantToken: nil,
			wantError: token.ErrInvalidInputSyntax,
		},
		{
			name:      "missing token",
			tokenID:   "H8pDtoIFzlzIlZ1fiH8pDtoFiJM7telsDIF525jt",
			wantToken: nil,
			wantError: token.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgTokenRepository(db)
			token, err := repo.FindByID(tt.tokenID)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if token != nil {
				if !reflect.DeepEqual(*token, *tt.wantToken) {
					t.Errorf("want \n%v\n; got \n%v\n", tt.wantToken, token)
				}
			}

		})
	}

}

// TestCreate tests token creation
func TestTokenSave(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name      string
		newToken  *token.Token
		wantError error
	}{
		{
			name: "Valid token",
			newToken: &token.Token{
				ID:        "ZumGB0WNQYKNVdSRbixKg32kkUqtKy69IXUaX275",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(30 * time.Minute).Round(time.Second),
				Kind:      "PwReset",
			},
			wantError: nil,
		},
		{
			name: "duplicate ID",
			newToken: &token.Token{
				ID:        "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(30 * time.Minute).Round(time.Second),
				Kind:      "PwReset",
			},
			wantError: token.ErrDuplicateField,
		},
		{
			name: "bad UID",
			newToken: &token.Token{
				ID:        "ZumGB0WNQYKNVdSRbixKg32kkUqtKy69IXUaX275",
				UID:       "dcce1beb-b724-94d470817323",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(30 * time.Minute).Round(time.Second),
				Kind:      "PwReset",
			},
			wantError: token.ErrInvalidInputSyntax,
		},
		{
			name: "empty ID",
			newToken: &token.Token{
				ID:        "",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(30 * time.Minute).Round(time.Second),
				Kind:      "PwReset",
			},
			wantError: token.ErrDb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgTokenRepository(db)

			_, err := repo.Save(tt.newToken)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
		})
	}
}

func TestTokenDelete(t *testing.T) {

	tests := []struct {
		name      string
		id        string
		wantError error
	}{
		{
			name:      "Existing token",
			id:        "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2", // see internal/storage/testdata/setup.sql
			wantError: nil,
		},
		{
			name:      "Non existing token",
			id:        "96db7333f9d549f8992903cbf1bdc24f", // random string
			wantError: token.ErrNoRecord,
		},
		{
			name:      "Empty token",
			id:        "",
			wantError: token.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgTokenRepository(db)
			err := repo.Delete(tt.id)
			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
		})
	}

}

func TestFindAllTokensByUID(t *testing.T) {
	tests := []struct {
		name      string
		UID       string
		tokensIDs *[]string
		wantError error
	}{
		{ // see pgk/storage/testdata/setup.sql
			name:      "User with tokens",
			UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
			tokensIDs: &[]string{"7telsDIFzlzIlZ1fiH8pDtoFiJMoBUi69j9525jt", "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2"},
			wantError: nil,
		},
		{
			name:      "User without tokens",
			UID:       "c835848e-88d1-4301-8927-2a83875d4b58",
			tokensIDs: nil,
			wantError: token.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgTokenRepository(db)
			tokens, err := repo.FindAllByUID(tt.UID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if (tokens != nil) &&
				(tt.tokensIDs != nil) &&
				len(*tokens) > 0 &&
				len(*tt.tokensIDs) > 0 {
				for _, id := range *tt.tokensIDs {
					if !containsToken(tokens, id) {
						t.Errorf("returned array does not contains %v", id)
					}
				}
			}
		})
	}
}

// Contains tells whether a contains x.
func containsToken(tokens *[]token.Token, id string) bool {
	for _, s := range *tokens {
		if s.ID == id {
			return true
		}
	}
	return false
}
