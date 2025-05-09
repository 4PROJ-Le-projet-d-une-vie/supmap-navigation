package incidents

import (
	"context"
	"encoding/json"
	"supmap-navigation/internal/gis"
	"supmap-navigation/internal/navigation"
	"supmap-navigation/internal/ws"
)

type Multicaster struct {
	Manager      *ws.Manager
	SessionCache navigation.SessionCache
}

func NewMulticaster(manager *ws.Manager, sessionCache navigation.SessionCache) *Multicaster {
	return &Multicaster{Manager: manager, SessionCache: sessionCache}
}

func (m *Multicaster) MulticastIncident(ctx context.Context, incident *Incident, action string) {
	m.Manager.RLock()
	defer m.Manager.RUnlock()
	for userID, client := range m.Manager.ClientsUnsafe() {
		session, err := m.SessionCache.GetSession(ctx, userID)
		if err != nil || session == nil {
			continue
		}
		if gis.IsPointInPolyline(gis.Point{Lat: incident.Lat, Lon: incident.Lon},
			convertNavPointsToGIS(session.Route.Polyline), 30) {
			payload, _ := json.Marshal(map[string]interface{}{
				"incident": incident,
				"action":   action,
			})
			client.Send(ws.Message{
				Type: "incident",
				Data: payload,
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
