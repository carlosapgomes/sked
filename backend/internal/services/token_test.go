package services_test

import (
	"testing"

	"carlosapgomes.com/gobackend/internal/mocks"
	"carlosapgomes.com/gobackend/internal/services"
	"carlosapgomes.com/gobackend/internal/token"
)

func TestTokenNew(t *testing.T) {
	repo := mocks.NewTokenRepo()
	svc := services.NewTokenService(repo)

	tests := []struct {
		name      string
		uid       string // uuidV4
		kind      string
		wantError error
	}{
		{
			name:      "Valid uid and kind",
			uid:       "6ba0a117-9eb8-412b-b6d4-52d096ff0e6a",
			kind:      token.PwReset,
			wantError: nil,
		},
		{
			name:      "Valid uid and kind2",
			uid:       "6ba0a117-9eb8-412b-b6d4-52d096ff0e6a",
			kind:      token.ValidateEmail,
			wantError: nil,
		},
		{
			name:      "Invalid uid",
			uid:       "6ba0a117-9eb8-0-52d096ff0e6a",
			kind:      token.PwReset,
			wantError: token.ErrInvalidInputSyntax,
		},
		{
			name:      "Invalid kind",
			uid:       "6ba0a117-9eb8-412b-b6d4-52d096ff0e6a",
			kind:      "my_special_kind",
			wantError: token.ErrInvalidInputSyntax,
		},
		{
			name:      "Mocked DB error",
			uid:       "44bd844b-a0b9-4d48-a371-cc2a44a7a4e3",
			kind:      token.PwReset,
			wantError: token.ErrDb,
		},
		{
			name:      "Wrong returned ID",
			uid:       "fbf344b4-2b23-4de6-9541-9678b904c020",
			kind:      token.PwReset,
			wantError: token.ErrDb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.New(tt.uid, tt.kind)
			if tt.wantError != err {
				t.Errorf("want error: %v, got %v", tt.wantError, err)
			}
			if id != nil {
				t.Logf("Received id = %s", *id)
			}
			if (id == nil) && (tt.wantError == nil) {
				t.Error("want valid token ID but got nill")
			}
		})
	}
}

func TestDeleteToken(t *testing.T) {
	repo := mocks.NewTokenRepo()
	svc := services.NewTokenService(repo)

	tests := []struct {
		name      string
		tID       string
		wantError error
	}{
		{
			name:      "Valid token ID",
			tID:       "iJiY6njKPdfPxZksvv27vM7C5cM8myaGQ4rnAzwm",
			wantError: nil,
		},
		{
			name:      "Non-existing token ID",
			tID:       "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2",
			wantError: token.ErrNoRecord,
		},
		{
			name:      "Bad token ID",
			tID:       "",
			wantError: token.ErrInvalidInputSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Delete(tt.tID)
			if tt.wantError != err {
				t.Errorf("want error: %v, got %v", tt.wantError, err)
			}
		})
	}
}
func TestFindAllByUIDToken(t *testing.T) {
	repo := mocks.NewTokenRepo()
	svc := services.NewTokenService(repo)
	tests := []struct {
		name       string
		uID        string
		wantTokens *[]string
		wantError  error
	}{
		{
			name:       "Valid id",
			uID:        "dcce1beb-aee6-4a4d-b724-94d470817323",
			wantTokens: &[]string{"7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2", "iJiY6njKPdfPxZksvv27vM7C5cM8myaGQ4rnAzwm"},
			wantError:  nil,
		},
		{
			name:       "Non-existing valid uid",
			uID:        "25fdc177-24ec-433e-955e-be5b240a4fb6",
			wantTokens: nil,
			wantError:  token.ErrNoRecord,
		},
		{
			name:       "Invalid uid",
			uID:        "dcce1beb-aee6-b724-94d470817323",
			wantTokens: nil,
			wantError:  token.ErrInvalidInputSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := svc.FindAllByUID(tt.uID)
			if (tokens != nil) &&
				(tt.wantTokens != nil) &&
				len(*tokens) > 0 &&
				len(*tt.wantTokens) > 0 {
				for _, id := range *tt.wantTokens {
					if !containsToken(tokens, id) {
						t.Errorf("returned array does not contains %v", id)
					}
				}
			}
			if tt.wantError != err {
				t.Errorf("want error: %v, got %v", tt.wantError, err)
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

func TestFindByIDToken(t *testing.T) {
	repo := mocks.NewTokenRepo()
	svc := services.NewTokenService(repo)

	tests := []struct {
		name      string
		tID       string
		wantError error
	}{
		{
			name:      "Valid id",
			tID:       "7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2",
			wantError: nil,
		},
		{
			name:      "Non-existing valid sid",
			tID:       "DCDNO4U5UC6I9632BndMkSFDCDNO4i1DCDNO4Gd2",
			wantError: token.ErrNoRecord,
		},
		{
			name:      "Invalid uid",
			tID:       "",
			wantError: token.ErrInvalidInputSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.FindByID(tt.tID)
			if tt.wantError != err {
				t.Errorf("want error: %s, got %s", tt.wantError, err)
			}
			if err == nil &&
				token != nil &&
				token.ID != tt.tID {
				t.Errorf("want ID: %s, got %s", tt.tID, token.ID)
			}
		})
	}
}
