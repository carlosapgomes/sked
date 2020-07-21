package session

// Repository interface for a session
type Repository interface {
	Save(s *Session) (*string, error)
	Get(sid string) (*Session, error)
	Delete(sid string) error
	FindAllByUID(uid string) (*[]Session, error)
}
