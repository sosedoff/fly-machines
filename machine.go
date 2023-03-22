package machines

import (
	"fmt"
)

type Machine struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	State      State    `json:"state"`
	Region     string   `json:"region"`
	InstanceID string   `json:"instance_id"`
	PrivateIP  string   `json:"private_ip"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	ImageRef   ImageRef `json:"image_ref"`
	Events     []Event  `json:"events"`
	Config     Config   `json:"config"`
}

type ImageRef struct {
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Digest     string `json:"digest"`
}

func (m Machine) CanStop() bool {
	switch m.State {
	case StateStarting, StateStarted:
		return true
	default:
		return false
	}
}

func (m Machine) CanDelete() bool {
	switch m.State {
	case StateStarting, StateStarted, StateStopping, StateStopped:
		return true
	default:
		return false
	}
}

func (m Machine) Inspect() string {
	return fmt.Sprintf(
		"machine(id=%q instance_id=%q region=%q state=%q created_at=%q updated_at=%q)",
		m.ID,
		m.InstanceID,
		m.Region,
		m.State,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
