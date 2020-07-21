package user

// Service interface for user model
type Service interface {
	Create(name, email, password, phone string) (*string, error)
	Authenticate(email, password string) (*string, error)
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	UpdateName(id, name string) error
	UpdateEmail(id, email string) error
	UpdatePhone(id, phone string) error
	UpdateRoles(id string, roles []string) error
	UpdatePw(id, password string) error
	UpdateStatus(id string, active bool) error
	UpdateEmailValidated(id string, active bool) error
	GetAll(before string, after string, pgSize int) (*Cursor, error)
	FindByName(name string, before string, hasBef bool, after string, hasAft bool, pgSize int) (*[]User, error)
}
