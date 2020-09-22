package mocks

import (
	"time"

	"carlosapgomes.com/sked/internal/session"
	uuid "github.com/satori/go.uuid"
)

type sessionDB map[string]session.Session

// SessionMockRepo is a session mock repository
type SessionMockRepo struct {
	db sessionDB
}

// NewSessionRepo returns a mock repository
func NewSessionRepo() *SessionMockRepo {
	db := make(map[string]session.Session)
	dt := time.Now().UTC()
	db["4e59d0bc-2fca-4aad-98d9-858709b66598"] = session.Session{
		ID:        "4e59d0bc-2fca-4aad-98d9-858709b66598",
		UID:       "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
		CreatedAt: dt,
		ExpiresAt: dt.Add(time.Duration(20) * time.Minute),
	}
	db["4e66a385-c7cd-47de-9e3b-cdfe26eecad4"] = session.Session{
		ID:        "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
		UID:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
		CreatedAt: dt,
		ExpiresAt: dt.Add(time.Duration(20) * time.Minute),
	}
	db["ef5de07d-4e78-4d8a-ab32-3529507a5072"] = session.Session{
		ID:        "ef5de07d-4e78-4d8a-ab32-3529507a5072",
		UID:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
		CreatedAt: dt.Add(time.Duration(-60) * time.Minute),
		ExpiresAt: dt.Add(time.Duration(-30) * time.Minute),
	}
	db["03e77847-ff41-4be6-9fc8-166da95097b9"] = session.Session{
		ID:        "03e77847-ff41-4be6-9fc8-166da95097b9",
		UID:       "ecadbb28-14e6-4560-8574-809c6c54b9cb",
		CreatedAt: dt,
		ExpiresAt: dt.Add(time.Duration(20) * time.Minute),
	}
	db["167ced64-af16-45d2-bb08-e35233c04ad1"] = session.Session{
		ID:        "167ced64-af16-45d2-bb08-e35233c04ad1",
		UID:       "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		CreatedAt: dt,
		ExpiresAt: dt.Add(time.Duration(20) * time.Minute),
	}
	return &SessionMockRepo{
		db,
	}
}

// Save a session and returns its id
func (r *SessionMockRepo) Save(s *session.Session) (*string, error) {
	if !isValidUUID(s.ID) || (!isValidUUID(s.UID)) {
		return nil, session.ErrInvalidInputSyntax
	}
	r.db[s.ID] = *s
	return &s.ID, nil
}

// Get a session
func (r *SessionMockRepo) Get(sid string) (*session.Session, error) {
	if !isValidUUID(sid) {
		return nil, session.ErrInvalidInputSyntax
	}
	if s, ok := r.db[sid]; ok {
		return &s, nil
	}
	return nil, session.ErrNoRecord

}

// Delete a session
func (r *SessionMockRepo) Delete(sid string) error {
	if sid == "4e59d0bc-2fca-4aad-98d9-858709b66598" {
		return nil
	}
	if !isValidUUID(sid) {
		return session.ErrInvalidInputSyntax
	}
	if _, ok := r.db[sid]; ok {
		delete(r.db, sid)
		return nil
	}
	return session.ErrNoRecord
}

// FindAllByUID finds all sessions connected to a user ID
func (r *SessionMockRepo) FindAllByUID(uid string) (*[]session.Session, error) {
	if !isValidUUID(uid) {
		return nil, session.ErrInvalidInputSyntax
	}
	var res []session.Session
	for _, s := range r.db {
		if s.UID == uid {
			res = append(res, s)
		}
	}
	return &res, nil
}

func isValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}

// SessionMockSvc mocks session services
type SessionMockSvc struct {
	repo session.Repository
}

// NewSessionSvc returns an instance of session services mock
func NewSessionSvc() *SessionMockSvc {
	repo := NewSessionRepo()
	return &SessionMockSvc{
		repo,
	}
}

// Create mocks session creation
func (s SessionMockSvc) Create(uid string) (*string, error) {
	id := "bfc6bbcd-d809-408c-beab-599eac369edc"
	return &id, nil
}

// Read mocks session read
func (s SessionMockSvc) Read(sid string) (*session.Session, error) {
	session, err := s.repo.Get(sid)
	if err != nil {
		return nil, err
	}
	return session, nil
	// if sid == "4e66a385-c7cd-47de-9e3b-cdfe26eecad4" {
	// 	dt := time.Now().UTC()
	// 	return &session.Session{
	// 		ID:        "4e66a385-c7cd-47de-9e3b-cdfe26eecad4",
	// 		UID:       "68b1d5e2-39dd-4713-8631-a08100383a0f",
	// 		CreatedAt: dt,
	// 		ExpiresAt: dt.Add(time.Duration(20) * time.Minute),
	// 	}, nil
	// }
	// if sid == "4e59d0bc-2fca-4aad-98d9-858709b66598" {
	// 	dt := time.Now().UTC()
	// 	return &session.Session{
	// 		ID:        "4e59d0bc-2fca-4aad-98d9-858709b66598",
	// 		UID:       "85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd",
	// 		CreatedAt: dt,
	// 		ExpiresAt: dt.Add(time.Duration(20) * time.Minute),
	// 	}, nil
	// }
	// return nil, session.ErrNoRecord
}

// Delete removes a session from DB
func (s SessionMockSvc) Delete(sid string) error {
	return s.repo.Delete(sid)
	// if (sid == "4e66a385-c7cd-47de-9e3b-cdfe26eecad4") ||
	// 	(sid == "4e59d0bc-2fca-4aad-98d9-858709b66598") {
	// 	return nil
	// }
	// return session.ErrNoRecord
}

// FindAllByUID returns all sessions associated with a user ID
func (s SessionMockSvc) FindAllByUID(uid string) (*[]session.Session, error) {
	return nil, nil
}

// IsValid checks a session expiration
func (s SessionMockSvc) IsValid(sid string) bool {
	return false
}
