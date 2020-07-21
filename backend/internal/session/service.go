package session

// Service interface for session model
type Service interface {
	Create(uid string) (*string, error)          // uid - user ID to be associated with
	Read(sid string) (*Session, error)           // sid - session ID
	Delete(sid string) error                     // sid - session ID
	FindAllByUID(uid string) (*[]Session, error) // uid - user ID associated with
	IsValid(sid string) bool                     // sid - session ID
}
