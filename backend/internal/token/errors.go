package token

import "errors"

// ErrNoRecord flags user not found
var ErrNoRecord = errors.New("token: no matching record found")

// ErrDuplicateField flags use of duplicate email
var ErrDuplicateField = errors.New("token: duplicate field")

// ErrInvalidInputSyntax flags text format errors (uuid for example)
var ErrInvalidInputSyntax = errors.New("token: invalid input syntax")

// ErrDb flags other errors
var ErrDb = errors.New("token: DB error")
