package storage

import (
	"database/sql"
	"fmt"
	"time"

	"carlosapgomes.com/gobackend/internal/token"
	"github.com/lib/pq"
)

type tokenRepository struct {
	DB *sql.DB
}

// NewPgTokenRepository creates a new token repo instance
func NewPgTokenRepository(db *sql.DB) token.Repository {
	return &tokenRepository{
		db,
	}
}

// FindByID finds a token by its ID
func (r *tokenRepository) FindByID(id string) (*token.Token, error) {
	if id == "" {
		return nil, fmt.Errorf("%w", token.ErrInvalidInputSyntax)
	}
	var t token.Token
	stmt := `SELECT id, userid, created_at, expires_at, kind from tokens where id = $1`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&t.ID, &t.UID, &t.CreatedAt, &t.ExpiresAt, &t.Kind)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w", token.ErrNoRecord)
	} else if err != nil {
		return nil, fmt.Errorf("%w\n %s", token.ErrDb, err.Error())
	}
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	t.CreatedAt = t.CreatedAt.In(loc)
	t.ExpiresAt = t.ExpiresAt.In(loc)
	return &t, err
}

// FindByUID finds all tokens	 belonging to a user
func (r *tokenRepository) FindAllByUID(uid string) (*[]token.Token, error) {
	var result []token.Token
	stmt := `SELECT * FROM tokens WHERE userid = $1`
	rows, err := r.DB.Query(stmt, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t token.Token
		if err := rows.Scan(&t.ID, &t.UID, &t.CreatedAt, &t.ExpiresAt, &t.Kind); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if len(result) == 0 {
		return nil, token.ErrNoRecord
	}
	return &result, nil
}

// Create generates a new token
func (r *tokenRepository) Save(t *token.Token) (*string, error) {
	stmt := `INSERT INTO tokens (id, userid, created_at, expires_at,kind) VALUES(
		$1, $2, $3, $4, $5) Returning id;`
	var id string
	err := r.DB.QueryRow(stmt, t.ID, t.UID, t.CreatedAt, t.ExpiresAt, t.Kind).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return nil, fmt.Errorf("%w\n %s %s", token.ErrDuplicateField, pqErr.Message, pqErr.Column)
		case "22P02":
			return nil, fmt.Errorf("%w\n %s %s", token.ErrInvalidInputSyntax, pqErr.Message, pqErr.Column)
		}
		return nil, fmt.Errorf("%w\n %s %s", token.ErrDb, pqErr.Message, pqErr.Column)
	}
	return &id, err
}

// Delete erases a token
func (r *tokenRepository) Delete(id string) error {
	stmt := `DELETE FROM tokens WHERE id = $1`
	res, err := r.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err == nil {
		if rows == 0 {
			return token.ErrNoRecord
		}
	}
	return nil
}
