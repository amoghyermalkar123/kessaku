package errors

import "errors"

var (
	ErrPoolReachedCapacity = errors.New("no worker is free and pool has reached capacity")
)
