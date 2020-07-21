package session

import "errors"

// ErrNoRecord flags user not found
var ErrNoRecord = errors.New("session: no matching record found")

// ErrDuplicateField flags use of duplicate email
var ErrDuplicateField = errors.New("session: duplicate field")

// ErrInvalidInputSyntax flags text format errors (uuid for example)
var ErrInvalidInputSyntax = errors.New("session: invalid input syntax")
