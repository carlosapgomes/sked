package storage_test

import (
	"reflect"
	"testing"
	"time"

	"carlosapgomes.com/sked/internal/session"
	"carlosapgomes.com/sked/internal/storage"
	"github.com/pkg/errors"
)

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name        string
		sessionID   string
		wantSession *session.Session
		wantError   error
	}{
		{name: "Valid sessionID",
			sessionID: "144f1223-70f4-4fca-9c99-c58eed7a0f4a",
			wantSession: &session.Session{
				ID:        "144f1223-70f4-4fca-9c99-c58eed7a0f4a",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Date(2019, 06, 23, 17, 0, 0, 0, time.UTC),
				ExpiresAt: time.Date(2019, 06, 23, 17, 0, 0, 0, time.UTC),
			},
			wantError: nil,
		},
		{name: "Invalid sessionID",
			sessionID: "144f1223-70f4-4fca",
			wantSession: &session.Session{
				ID:        "144f1223-70f4-4fca-9c99-c58eed7a0f4a",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Date(2019, 06, 23, 17, 0, 0, 0, time.UTC),
				ExpiresAt: time.Date(2019, 06, 23, 17, 0, 0, 0, time.UTC),
			},
			wantError: session.ErrInvalidInputSyntax,
		},
		{name: "Missing session",
			sessionID:   "5c5eb8f0-ea75-4d2f-bde8-75f91f817256",
			wantSession: nil,
			wantError:   session.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			repo := storage.NewPgSessionRepository(db)
			session, err := repo.Get(tt.sessionID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if session != nil {
				if !reflect.DeepEqual(*session, *tt.wantSession) {
					t.Errorf("want \n%v\n; got \n%v\n", tt.wantSession, session)
				}
			}

		})
	}
}

func TestSave(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name       string
		newSession *session.Session
		wantError  error
	}{
		{
			name: "Valid Session",
			newSession: &session.Session{
				ID:        "9df27193-4e40-43e7-8778-60a4688ced3f",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(20 * time.Minute).Round(time.Second),
			},
			wantError: nil,
		},
		{
			name: "Duplicate Session",
			newSession: &session.Session{
				ID:        "144f1223-70f4-4fca-9c99-c58eed7a0f4a",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(20 * time.Minute).Round(time.Second),
			},
			wantError: session.ErrDuplicateField,
		},
		{
			name: "Invalid ID",
			newSession: &session.Session{
				ID:        "9df27193",
				UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(20 * time.Minute).Round(time.Second),
			},
			wantError: session.ErrInvalidInputSyntax,
		},
		{
			name: "Invalid UID",
			newSession: &session.Session{
				ID:        "9df27193-4e40-43e7-8778-60a4688ced3f",
				UID:       "dcce1beb-",
				CreatedAt: time.Now().UTC().Round(time.Second),
				ExpiresAt: time.Now().UTC().Add(20 * time.Minute).Round(time.Second),
			},
			wantError: session.ErrInvalidInputSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgSessionRepository(db)

			id, err := repo.Save(tt.newSession)
			if err == nil {
				if err != tt.wantError {
					t.Errorf("want %v; got %s", tt.wantError, err)
				}
			} else {
				if errors.Cause(err) != tt.wantError {
					t.Errorf("want %v; got %s", tt.wantError, errors.Cause(err))
				}
			}
			if id != nil {
				session, err := repo.Get(*id)
				if err != nil {
					t.Errorf("want \n%v\n; got %s", tt.newSession, err)
				}
				if session != nil {
					if !reflect.DeepEqual(*session, *tt.newSession) {
						t.Errorf("want \n%v\n; got \n%v\n", tt.newSession, session)
					}
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {

	tests := []struct {
		name      string
		sid       string
		wantError error
	}{
		{
			name:      "Existing session",
			sid:       "144f1223-70f4-4fca-9c99-c58eed7a0f4a", // see internal/storage/testdata/setup.sql
			wantError: nil,
		},
		{
			name:      "Non existing session",
			sid:       "96db7333-f9d5-49f8-9929-03cbf1bdc24f", // random uuidV4
			wantError: session.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgSessionRepository(db)
			err := repo.Delete(tt.sid)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
		})
	}

}

func TestFindAllByUID(t *testing.T) {
	tests := []struct {
		name      string
		UID       string
		sessions  *[]string
		wantError error
	}{
		{ // see pgk/storage/testdata/setup.sql
			name:      "User with sessions",
			UID:       "dcce1beb-aee6-4a4d-b724-94d470817323",
			sessions:  &[]string{"144f1223-70f4-4fca-9c99-c58eed7a0f4a", "ddfd72b2-d57b-4fc4-9265-55b9dae7287c"},
			wantError: nil,
		},
		{
			name:      "User without sessions",
			UID:       "c835848e-88d1-4301-8927-2a83875d4b58",
			sessions:  nil,
			wantError: session.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			repo := storage.NewPgSessionRepository(db)
			sessions, err := repo.FindAllByUID(tt.UID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if (sessions != nil) &&
				len(*sessions) > 0 &&
				(tt.sessions != nil) &&
				len(*tt.sessions) > 0 {
				for _, sid := range *tt.sessions {
					if !contains(sessions, sid) {
						t.Errorf("returned array does not contains %v", sid)
					}
				}
			}
		})
	}
}

// Contains tells whether a contains x.
func contains(sessions *[]session.Session, sid string) bool {
	for _, s := range *sessions {
		if s.ID == sid {
			return true
		}
	}
	return false
}
