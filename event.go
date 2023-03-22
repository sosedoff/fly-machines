package machines

import (
	"fmt"
	"time"
)

type Event struct {
	ID        string        `json:"id"`
	Type      string        `json:"type"`
	Status    string        `json:"status"`
	Request   *EventRequest `json:"request"`
	Source    string        `json:"source"`
	Timestamp int64         `json:"timestamp"`
}

type EventRequest struct {
	ExitEvent *ExitEvent `json:"exit_event,omitempty"`
}

type ExitEvent struct {
	ExitCode      int       `json:"exit_code,omitempty"`
	ExitedAt      time.Time `json:"exited_at,omitempty"`
	GuestExitCode int       `json:"guest_exit_code,omitempty"`
	GuestSignal   int       `json:"guest_signal,omitempty"`
	OOMKilled     bool      `json:"oom_killed,omitempty"`
	RequestedStop bool      `json:"requested_stop,omitempty"`
	Restarting    bool      `json:"restarting,omitempty"`
	Signal        int       `json:"signal,omitempty"`
}

func (e Event) Inspect() string {
	return fmt.Sprintf("event(id=%q type=%q status=%q source=%q ts=%d)",
		e.ID,
		e.Type,
		e.Status,
		e.Source,
		e.Timestamp,
	)
}
