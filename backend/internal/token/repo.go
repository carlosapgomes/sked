package token

// Repository interface for token model
type Repository interface {
	FindByID(string) (*Token, error)
	FindAllByUID(string) (*[]Token, error)
	Save(token *Token) (*string, error)
	Delete(string) error
}
