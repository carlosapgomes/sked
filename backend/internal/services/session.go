package services

import (
	"time"

	"carlosapgomes.com/sked/internal/session"
	uuid "github.com/satori/go.uuid"
)

type sessionService struct {
	sessionLifetime int // session lifetime in minutes
	repo            session.Repository
}

// NewSessionService returns a session Service instance
func NewSessionService(sessionLifetime int,
	repo session.Repository) session.Service {
	return &sessionService{
		sessionLifetime,
		repo,
	}
}

// Create, creates a new session and returns a session ID
func (s *sessionService) Create(uid string) (*string, error) {
	dt := time.Now().UTC()
	session := &session.Session{
		ID:        uuid.NewV4().String(),
		UID:       uid,
		CreatedAt: dt,
		ExpiresAt: dt.Add(time.Duration(s.sessionLifetime) * time.Minute),
	}
	sid, err := s.repo.Save(session)
	if err != nil {
		return nil, err
	}
	if *sid != session.ID {
		return nil, nil
	}
	return &session.ID, nil
}

func (s *sessionService) Read(sid string) (*session.Session, error) {
	session, err := s.repo.Get(sid)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *sessionService) Delete(sid string) error {
	return s.repo.Delete(sid)
}

func (s *sessionService) FindAllByUID(uid string) (*[]session.Session, error) {
	sessions, err := s.repo.FindAllByUID(uid)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *sessionService) IsValid(sid string) bool {
	session, err := s.Read(sid)
	if err != nil {
		return false
	}
	return (time.Now().UTC().Before(session.ExpiresAt))
}
