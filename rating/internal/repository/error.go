package repository

import "errors"

// ErrorNotFound returns when a request record is not found.
var ErrNotFound = errors.New("not found")
