package machines

import (
	"errors"
)

var (
	ErrAppNameRequired   = errors.New("app name is required")
	ErrAuthRequired      = errors.New("api token is required")
	ErrInvalidAuth       = errors.New("invalid or expired auth token")
	ErrInputRequired     = errors.New("request input required")
	ErrMachineIDRequired = errors.New("machine id is required")
	ErrInvalidWaitState  = errors.New("state must be one of started/stopped/destroyed")
)
