// user repository port interface definition

package user

// Repository inteface definition for user model
type Repository interface {
	Create(user User) (*string, error)
	UpdatePw(id string, password []byte) error
	UpdateStatus(id string, active bool) error
	UpdateName(id, name string) error
	UpdateEmail(id, email string) error
	UpdatePhone(id, phone string) error
	UpdateEmailValidated(id string, active bool) error
	UpdateRoles(id string, roles []string) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	GetAll(cursor string, after bool, pgSize int) (*[]User, bool, error)
	GetAllDoctors() (*[]User, error)
	FindByName(name string) (*[]User, error)
}
