package patient

// Service interface for user model
type Service interface {
	Create(name, address, city, state string, phone []string, createdBy string) (*string, error)
	FindByID(id string) (*Patient, error)
	UpdateName(id, name string) error
	UpdatePhone(id string, phones []string) error
	GetAll(before string, after string, pgSize int) (*Page, error)
	FindByName(name string) (*[]Patient, error)
}
