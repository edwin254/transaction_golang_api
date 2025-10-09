package errors

import "errors"

var (
    ErrNotFound          = errors.New("transaction not found")
    ErrInvalidTransition = errors.New("invalid status transition")
)
