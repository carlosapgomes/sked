// user repository port interface definition

package patient

// Repository inteface definition for user model
type Repository interface {
	Create(patient Patient) (*string, error)
	UpdateName(id, name string) error
	UpdatePhone(id, phone string) error
	FindByID(id string) (*Patient, error)
	GetAll(cursor string, after bool, pgSize int) (*[]Patient, bool, error)
	FindByName(name string, before string, hasBef bool, after string, hasAft bool, pgSize int) (*[]Patient, error)
}
