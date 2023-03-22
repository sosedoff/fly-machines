package machines

import (
	"time"
)

type ListInput struct {
	State State
}

type GetInput struct {
	ID string
}

type CreateInput struct {
	Name   string            `json:"name,omitempty"`
	Region string            `json:"region,omitempty"`
	Config *Config           `json:"config"`
	Size   Size              `json:"size,omitempty"`
	Env    map[string]string `json:"env,omitempty"`
}

type StopInput struct {
	ID      string        `json:"id,omitempty"`
	Signal  int           `json:"signal,omitempty"`
	Timeout time.Duration `json:"timeout,omitempty"`
}

func (i StopInput) Validate() error {
	if i.ID == "" {
		return ErrMachineIDRequired
	}
	return nil
}

type DeleteInput struct {
	ID      string
	AppName string
	Kill    bool
}

func (i DeleteInput) Validate() error {
	if i.ID == "" {
		return ErrMachineIDRequired
	}
	return nil
}

type WaitInput struct {
	ID         string
	AppName    string
	InstanceID string
	State      State
	Timeout    time.Duration
}

func (i WaitInput) Validate() error {
	if i.ID == "" {
		return ErrMachineIDRequired
	}

	switch i.State {
	case StateStarted, StateStopped, StateDestroyed:
		return nil
	default:
		return ErrInvalidWaitState
	}
}

type LeaseInput struct {
	ID    string `json:"-"`
	Nonce string `json:"-"`
	TTL   int    `json:"ttl,omitempty"`
}

func (i LeaseInput) Validate() error {
	if i.ID == "" {
		return ErrMachineIDRequired
	}
	return nil
}
