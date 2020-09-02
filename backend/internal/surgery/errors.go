package surgery

import "errors"

// ErrNoRecord flags surgery not found
var ErrNoRecord = errors.New("surgeries: no matching record found")

// ErrDuplicateField flags use of duplicate
var ErrDuplicateField = errors.New("surgeries: duplicate field")

// ErrInvalidInputSyntax flags text format errors (uuid for example)
var ErrInvalidInputSyntax = errors.New("surgeries: invalid input syntax")

// ErrDb flags other errors
var ErrDb = errors.New("surgeries: DB error")
