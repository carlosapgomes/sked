package token

// Service interface for token model
type Service interface {
	New(uid string, kind string) (*string, error)
	FindByID(id string) (*Token, error)
	FindAllByUID(uid string) (*[]Token, error)
	Delete(id string) error
}
