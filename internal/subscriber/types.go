package subscriber

import "supmap-navigation/internal/incidents"

// IncidentMessage represents any message received in the incidents pub/sub channel.
type IncidentMessage struct {
	Data   incidents.Incident `json:"data"`
	Action Action             `json:"action"`
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
