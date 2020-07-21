package user

import "errors"

// ErrNoRecord flags user not found
var ErrNoRecord = errors.New("users: no matching record found")

// ErrInvalidCredentials flags unsuccesfull login
var ErrInvalidCredentials = errors.New("users: invalid credentials")

// ErrDuplicateEmail flags use of duplicate email
var ErrDuplicateEmail = errors.New("users: duplicate email")

// ErrDuplicateField flags use of duplicate email
var ErrDuplicateField = errors.New("users: duplicate field")

// ErrInvalidInputSyntax flags text format errors (uuid for example)
var ErrInvalidInputSyntax = errors.New("users: invalid input syntax")

// ErrEmailIsNotValidated flags text format errors (uuid for example)
var ErrEmailIsNotValidated = errors.New("users: email address is not validated")

// ErrDb flags other errors
var ErrDb = errors.New("users: DB error")
