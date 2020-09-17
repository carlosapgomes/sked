package storage

import (
	"database/sql"
	"fmt"
	"time"

	"carlosapgomes.com/sked/internal/user"
	"github.com/lib/pq"
)

// userRepository type
type userRepository struct {
	DB *sql.DB
}

// NewPgUserRepository returns an instance of a user repo
func NewPgUserRepository(db *sql.DB) user.Repository {
	return &userRepository{
		db,
	}
}

// CreateUser creates a new user
func (r userRepository) Create(u user.User) (*string, error) {
	stmt := `INSERT INTO users (id, name, email, phone, hashedpw, created_at, updated_at,
 Roles ) VALUES($1, $2, $3, $4, $5, $6, $7, $8) Returning id;`

	var id string
	err := r.DB.QueryRow(stmt, u.ID, u.Name, u.Email, u.Phone, u.HashedPw, u.CreatedAt, u.UpdatedAt, pq.Array(u.Roles)).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return nil, fmt.Errorf("%w\n %s %s", user.ErrDuplicateField, pqErr.Message, pqErr.Column)
		case "22P02":
			return nil, fmt.Errorf("%w\n %s %s", user.ErrInvalidInputSyntax, pqErr.Message, pqErr.Column)
		}
		return nil, fmt.Errorf("%w\n %s %s", user.ErrDb, pqErr.Message, pqErr.Column)
	}
	return &id, err
}

// UpdateUserPw updates user password
func (r userRepository) UpdatePw(id string, pw []byte) error {
	stmt := `UPDATE users SET hashedpw = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(stmt, pw, time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// UpdateStatus updates user account status (active = true|false)
func (r userRepository) UpdateStatus(id string, active bool) error {
	stmt := `UPDATE users SET active = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(stmt, active, time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// UpdateEmailValidated updates user account status (active = true|false)
func (r userRepository) UpdateEmailValidated(id string, validated bool) error {
	stmt := `UPDATE users SET email_was_validated = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(stmt, validated, time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// UpdateName updates user name
func (r userRepository) UpdateName(id string, name string) error {
	stmt := `UPDATE users SET "name" = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(stmt, name, time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// UpdatePhone updates user phone number
func (r userRepository) UpdatePhone(id string, phone string) error {
	stmt := `UPDATE users SET "phone" = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(stmt, phone, time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// UpdateEmail updates user email
func (r userRepository) UpdateEmail(id string, email string) error {
	stmt := `UPDATE users SET email = $1, updated_at = $2, email_was_validated = $3 WHERE id = $4`
	_, err := r.DB.Exec(stmt, email, time.Now().UTC(), false, id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err
}

// UpdateRoles set user roles, it replaces the older array completely
func (r userRepository) UpdateRoles(id string, roles []string) error {
	stmt := `UPDATE users SET Roles = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(stmt, pq.Array(roles), time.Now().UTC(), id)
	if err != nil {
		pqErr := err.(*pq.Error)
		return pqErr
	}
	return err

}

// FindByID searchs a user by its id
func (r userRepository) FindByID(id string) (*user.User, error) {
	var u user.User
	stmt := `SELECT name, email, phone, hashedpw, created_at, updated_at, active, email_was_validated, roles
	from users where id = $1`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&u.Name, &u.Email, &u.Phone, &u.HashedPw, &u.CreatedAt, &u.UpdatedAt,
		&u.Active, &u.EmailWasValidated, pq.Array(&u.Roles))
	if err == sql.ErrNoRows {
		return nil, user.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	u.ID = id
	// hashed user password will be return to the calling service
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	u.CreatedAt = u.CreatedAt.In(loc)
	u.UpdatedAt = u.UpdatedAt.In(loc)
	return &u, err
}

// FindByEmail searchs a user by its id
func (r userRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	stmt := `SELECT id, name, email, phone, hashedpw, created_at, updated_at, active, email_was_validated, roles
	from users where email = $1`
	row := r.DB.QueryRow(stmt, email)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.HashedPw, &u.CreatedAt, &u.UpdatedAt,
		&u.Active, &u.EmailWasValidated, pq.Array(&u.Roles))
	// hashed user password will be return to the calling service
	// u.HashedPw = []byte{}
	if err == sql.ErrNoRows {
		return nil, user.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	u.CreatedAt = u.CreatedAt.In(loc)
	u.UpdatedAt = u.UpdatedAt.In(loc)
	return &u, err
}

// GetAll returns a paginated list of all users ordered by email and a bool if there are more results in this direction
func (r userRepository) GetAll(cursor string, next bool, pgSize int) (*[]user.User, bool, error) {
	if pgSize <= 0 {
		pgSize = 15
	}
	var stmt string
	if next {
		// Get next result page
		stmt = `SELECT id, name, email, phone, created_at, updated_at, active, email_was_validated, roles
		FROM users WHERE email > $1 ORDER BY email LIMIT $2`
	} else {
		// Get previous result page
		stmt = `SELECT id, name, email, phone, created_at, updated_at, active, email_was_validated, roles
		FROM users WHERE email < $1 ORDER BY email LIMIT $2`
	}
	rows, err := r.DB.Query(stmt, cursor, (pgSize + 1))
	if err != nil {
		return nil, false, err
	}
	var users []user.User
	defer rows.Close()
	for rows.Next() {
		var u user.User
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.CreatedAt, &u.UpdatedAt,
			&u.Active, &u.EmailWasValidated, pq.Array(&u.Roles))
		if err == sql.ErrNoRows {
			return nil, false, user.ErrNoRecord
		} else if err != nil {
			return nil, false, err
		}
		users = append(users, u)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, false, err
	}
	hasMore := false
	if len(users) == pgSize {
		// remove last slice item, because it was not requested
		users = users[:len(users)-1]
		hasMore = true
	}
	return &users, hasMore, nil
}

// FindByName returns a list of users whose names looks like 'name'
func (r userRepository) FindByName(name string) (*[]user.User, error) {
	return nil, nil
}
