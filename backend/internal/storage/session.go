package storage

import (
	"database/sql"
	"time"

	"carlosapgomes.com/gobackend/internal/session"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type sessionRepository struct {
	DB *sql.DB
}

// NewPgSessionRepository creates a new session repo
func NewPgSessionRepository(db *sql.DB) session.Repository {
	return &sessionRepository{
		db,
	}
}

// Save, saves a session
func (r *sessionRepository) Save(s *session.Session) (*string, error) {
	stmt := `INSERT INTO sessions (id, userid, created_at, expires_at) VALUES(
		$1, $2, $3, $4) Returning id;`
	var id string
	err := r.DB.QueryRow(stmt, s.ID, s.UID, s.CreatedAt, s.ExpiresAt).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return nil, errors.Wrapf(session.ErrDuplicateField, "%s %s", pqErr.Message, pqErr.Column)
		case "22P02":
			return nil, errors.Wrapf(session.ErrInvalidInputSyntax, "%s %s", pqErr.Message, pqErr.Column)
		default:
			return nil, errors.Wrapf(pqErr, "code %s %s %s", pqErr.Code, pqErr.Message, pqErr.Column)
		}
	}
	return &id, err
}

// Get, returns a session
func (r *sessionRepository) Get(sid string) (*session.Session, error) {
	sessionID, err := uuid.FromString(sid)
	if err != nil {
		return nil, session.ErrInvalidInputSyntax
	}
	var s session.Session
	stmt := `SELECT id, userid, created_at, expires_at from sessions where id = $1`
	row := r.DB.QueryRow(stmt, sessionID)
	err = row.Scan(&s.ID, &s.UID, &s.CreatedAt, &s.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, session.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// every date/time was saved as UTC, so use them as UTC
	loc, _ := time.LoadLocation("UTC")
	s.CreatedAt = s.CreatedAt.In(loc)
	s.ExpiresAt = s.ExpiresAt.In(loc)
	return &s, err
}

// Delete, destroy a session
func (r *sessionRepository) Delete(sid string) error {
	stmt := `DELETE FROM sessions WHERE id = $1`
	res, err := r.DB.Exec(stmt, sid)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err == nil {
		if rows == 0 {
			return session.ErrNoRecord
		}
	}
	return nil
}

// FindAllByUID returns an array of all sessions associated with a user
func (r *sessionRepository) FindAllByUID(uid string) (*[]session.Session, error) {
	var result []session.Session
	stmt := `SELECT * FROM sessions WHERE userid = $1`
	rows, err := r.DB.Query(stmt, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s session.Session
		if err := rows.Scan(&s.ID, &s.UID, &s.CreatedAt, &s.ExpiresAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if len(result) == 0 {
		return nil, session.ErrNoRecord
	}
	return &result, nil
}
