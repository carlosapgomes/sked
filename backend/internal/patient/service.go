package patient

// Service interface for user model
type Service interface {
	Create(name, address, city, state string,
		phone []string, createdBy string) (*string, error)
	FindByID(id string) (*Patient, error)
	UpdateName(id, name, updatedBy string) error
	UpdatePhone(id string, phones []string, updatedBy string) error
	GetAll(before string, after string, pgSize int) (*Page, error)
	UpdatePatient(id, name, address, city, state string,
		phone []string, updatedBy string) error
	FindByName(name string) (*[]Patient, error)
}
