package patient

import "errors"

// ErrNoRecord flags user not found
var ErrNoRecord = errors.New("patients: no matching record found")

// ErrDuplicateField flags use of duplicate email
var ErrDuplicateField = errors.New("patients: duplicate field")

// ErrInvalidInputSyntax flags text format errors (uuid for example)
var ErrInvalidInputSyntax = errors.New("patients: invalid input syntax")

// ErrDb flags other errors
var ErrDb = errors.New("patients: DB error")
