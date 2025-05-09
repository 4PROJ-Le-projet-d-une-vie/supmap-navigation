package incidents

import (
	"time"
)

// IncidentPayload represents the payload used in the messages sent to the clients.
type IncidentPayload struct {
	Incident *Incident `json:"incident"`
	Action   string    `json:"action"`
}

type Action string

const (
	Create    Action = "create"
	Certified Action = "certified"
	Deleted   Action = "deleted"
)

func (a *Action) IsValid() bool {
	switch *a {
	case Create, Certified, Deleted:
		return true
	}
	return false
}

type Incident struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Type      *Type      `json:"type"`
	Lat       float64    `json:"lat"`
	Lon       float64    `json:"lon"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Type struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	NeedRecalculation bool   `json:"need_recalculation"`
}
