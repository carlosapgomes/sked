package appointment

import "errors"

// ErrNoRecord flags user not found
var ErrNoRecord = errors.New("appointments: no matching record found")

// ErrDuplicateField flags use of duplicate
var ErrDuplicateField = errors.New("appointments: duplicate field")

// ErrInvalidInputSyntax flags text format errors (uuid for example)
var ErrInvalidInputSyntax = errors.New("appointments: invalid input syntax")

// ErrDb flags other errors
var ErrDb = errors.New("appointments: DB error")
