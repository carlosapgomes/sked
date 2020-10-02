// user repository port interface definition

package patient

// Repository inteface definition for user model
type Repository interface {
	Create(patient Patient) (*string, error)
	UpdateName(id, name string) error
	UpdatePhone(id string, phones []string) error
	FindByID(id string) (*Patient, error)
	FindByName(name string) (*[]Patient, error)
	GetAll(cursor string, after bool, pgSize int) (*[]Patient, bool, error)
}
