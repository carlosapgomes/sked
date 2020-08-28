package services_test

import (
	"testing"

	"carlosapgomes.com/sked/internal/mocks"
	"carlosapgomes.com/sked/internal/services"
	"carlosapgomes.com/sked/internal/session"
	uuid "github.com/satori/go.uuid"
)

func TestSessionCreate(t *testing.T) {
	repo := mocks.NewSessionRepo()
	svc := services.NewSessionService(20, repo)

	tests := []struct {
		name      string
		uid       string // uuidV4
		wantError error
	}{
		{
			name:      "Valid uid",
			uid:       "6ba0a117-9eb8-412b-b6d4-52d096ff0e6a",
			wantError: nil,
		},
		{
			name:      "Invalid uid",
			uid:       "6ba0a117-9eb8-0-52d096ff0e6a",
			wantError: session.ErrInvalidInputSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.Create(tt.uid)
			if tt.wantError != err {
				t.Errorf("want error: %s, got %s", tt.wantError, err)
			}
			if (id == nil) && (tt.wantError == nil) {
				t.Errorf("want ID: %s, got nill", tt.uid)
			}
		})
	}
}

func TestReadSession(t *testing.T) {
	repo := mocks.NewSessionRepo()
	svc := services.NewSessionService(20, repo)

	uid := uuid.NewV4().String()
	sessionID, err := svc.Create(uid)
	if err != nil {
		t.Errorf("Could not creat new session, %s", err)
		return
	}
	tests := []struct {
		name      string
		sid       string // uuidV4
		wantError error
	}{
		{
			name:      "Valid sid",
			sid:       *sessionID,
			wantError: nil,
		},
		{
			name:      "Non-existing valid sid",
			sid:       "4592d639-0112-465d-a59d-2f93f3608b7a",
			wantError: session.ErrNoRecord,
		},
		{
			name:      "Invalid uid",
			sid:       "6ba0a117-9eb8-0-52d096ff0e6a",
			wantError: session.ErrInvalidInputSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := svc.Read(tt.sid)
			if tt.wantError != err {
				t.Errorf("want error: %s, got %s", tt.wantError, err)
			}
			if err == nil &&
				session != nil &&
				session.ID != tt.sid {
				t.Errorf("want ID: %s, got %s", tt.sid, session.ID)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	repo := mocks.NewSessionRepo()
	svc := services.NewSessionService(20, repo)

	uid := uuid.NewV4().String()
	sessionID, err := svc.Create(uid)
	if err != nil {
		t.Errorf("Could not creat new session, %s", err)
		return
	}
	tests := []struct {
		name      string
		sid       string // uuidV4
		wantError error
	}{
		{
			name:      "Valid sid",
			sid:       *sessionID,
			wantError: nil,
		},
		{
			name:      "Non-existing valid sid",
			sid:       "4592d639-0112-465d-a59d-2f93f3608b7a",
			wantError: session.ErrNoRecord,
		},
		{
			name:      "Invalid sid",
			sid:       "6ba0a117-9eb8-0-52d096ff0e6a",
			wantError: session.ErrInvalidInputSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Delete(tt.sid)
			if tt.wantError != err {
				t.Errorf("want error: %s, got %s", tt.wantError, err)
			}
		})
	}
}

func TestSessionFindAllByUID(t *testing.T) {
	repo := mocks.NewSessionRepo()
	svc := services.NewSessionService(20, repo)

	uid := uuid.NewV4().String()
	// create 03 sessions for user uid
	var sessions []string
	for i := 0; i < 3; i++ {
		sessionID, err := svc.Create(uid)
		if err != nil {
			t.Errorf("Could not creat new session, %s", err)
			return
		} else {
			sessions = append(sessions, *sessionID)
		}
	}

	tests := []struct {
		name            string
		uid             string // uuidV4
		wantError       error
		wantNOfSessions int
	}{
		{
			name:            "Valid uid with associated sessions",
			uid:             uid,
			wantError:       nil,
			wantNOfSessions: 3,
		},
		{
			name:            "Valid uid without associated sessions",
			uid:             "b0bca004-28eb-4289-984d-1a357592eafb",
			wantError:       nil,
			wantNOfSessions: 0,
		},
		{
			name:            "Invalid uid",
			uid:             "6ba0a117-9eb8-0-52d096ff0e6a",
			wantError:       session.ErrInvalidInputSyntax,
			wantNOfSessions: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessions, err := svc.FindAllByUID(tt.uid)
			if tt.wantError != err {
				t.Errorf("want error: %s, got %s", tt.wantError, err)
			}
			if tt.wantError == nil && sessions == nil {
				t.Errorf("want sessions not to be nil, got %v", sessions)
			}
			if (sessions != nil) && (len(*sessions) != tt.wantNOfSessions) {
				t.Errorf("Want %d sessions, got %d", tt.wantNOfSessions, len(*sessions))
			}

		})
	}
}

func TestIsValidSession(t *testing.T) {
	repo := mocks.NewSessionRepo()
	svc := services.NewSessionService(20, repo)

	// uid := uuid.NewV4().String()
	// sessionID, err := svc.Create(uid)
	// if err != nil {
	// 	t.Errorf("Could not creat new session, %s", err)
	// 	return
	// }
	tests := []struct {
		name    string
		sid     string // uuidV4
		wantRes bool
	}{
		{
			name:    "Non-expired session",
			sid:     "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
			wantRes: true,
		},
		{
			name:    "Expired session",
			sid:     "ef5de07d-4e78-4d8a-ab32-3529507a5072",
			wantRes: false,
		},
		{
			name:    "Non-existing session",
			sid:     "4e78e07d-ab32-4d8a-ab32-3529507a5072",
			wantRes: false,
		},
		// {
		// 	name:      "Non-existing valid sid",
		// 	sid:       "4592d639-0112-465d-a59d-2f93f3608b7a",
		// 	wantError: session.ErrNoRecord,
		// },
		// {
		// 	name:      "Invalid sid",
		// 	sid:       "6ba0a117-9eb8-0-52d096ff0e6a",
		// 	wantError: session.ErrInvalidInputSyntax,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := svc.IsValid(tt.sid)
			if tt.wantRes != res {
				t.Errorf("want error: %v, got %v", tt.wantRes, res)
			}
		})
	}
}
