package subscriber

import "supmap-navigation/internal/incidents"

// IncidentMessage represents any message received in the incidents pub/sub channel.
type IncidentMessage struct {
	Data   incidents.Incident `json:"data"`
	Action incidents.Action   `json:"action"`
}
