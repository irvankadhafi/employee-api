package usecase

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrDuplicateEmployee = errors.New("employee already exist")
	ErrPermissionDenied  = errors.New("permission denied")
)
