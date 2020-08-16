package patient

// Service interface for user model
type Service interface {
	Create(name, email, password, phone string) (*string, error)
	FindByID(id string) (*Patient, error)
	UpdateName(id, name string) error
	UpdatePhone(id, phone string) error
	GetAll(before string, after string, pgSize int) (*Cursor, error)
	FindByName(name string, before string, hasBef bool, after string, hasAft bool, pgSize int) (*[]Patient, error)
}
