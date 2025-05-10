package incidents

import (
	"context"
	"encoding/json"
	"supmap-navigation/internal/gis"
	routing "supmap-navigation/internal/gis/routing"
	"supmap-navigation/internal/navigation"
	"supmap-navigation/internal/ws"
)

type Multicaster struct {
	Manager       *ws.Manager
	SessionCache  navigation.SessionCache
	RoutingClient *routing.Client
}

func NewMulticaster(manager *ws.Manager, sessionCache navigation.SessionCache, routingClient *routing.Client) *Multicaster {
	return &Multicaster{
		Manager:       manager,
		SessionCache:  sessionCache,
		RoutingClient: routingClient,
	}
}

func (m *Multicaster) MulticastIncident(ctx context.Context, incident *Incident, action string) {
	m.Manager.RLock()
	defer m.Manager.RUnlock()
	for sessionID, client := range m.Manager.ClientsUnsafe() {
		session, err := m.SessionCache.GetSession(ctx, sessionID)
		if err != nil || session == nil {
			continue
		}
		if gis.IsPointInPolyline(gis.Point{Lat: incident.Lat, Lon: incident.Lon},
			convertNavPointsToGIS(session.Route.Polyline), 30) {

			incidentPayload := IncidentPayload{
				Incident: incident,
				Action:   action,
			}

			jsonPayload, _ := json.Marshal(incidentPayload)
			client.Send(ws.Message{
				Type: "incident",
				Data: jsonPayload,
			})
		}
	}
}

func convertNavPointsToGIS(points []navigation.Point) []gis.Point {
	res := make([]gis.Point, len(points))
	for i, p := range points {
		res[i] = gis.Point{Lat: p.Lat, Lon: p.Lon}
	}
	return res
}
