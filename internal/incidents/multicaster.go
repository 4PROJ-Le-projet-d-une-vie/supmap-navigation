package incidents

import (
	"context"
	"encoding/json"
	"log"
	"supmap-navigation/internal/gis"
	routing "supmap-navigation/internal/gis/routing"
	"supmap-navigation/internal/navigation"
	"supmap-navigation/internal/ws"
	"time"
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

// MulticastIncident notifies each client if it's impacted by the incident.
// If the incident needs a route recalculation and is certified, the new route is sent to the clients.
func (m *Multicaster) MulticastIncident(ctx context.Context, incident *Incident, action string) {
	m.Manager.RLock()
	defer m.Manager.RUnlock()

	for sessionID, client := range m.Manager.ClientsUnsafe() {
		session, err := m.SessionCache.GetSession(ctx, sessionID)
		if err != nil || session == nil {
			continue
		}
		if !m.isIncidentOnRoute(incident, session) {
			continue
		}

		if incident.Type != nil && action == "certified" && incident.Type.NeedRecalculation {
			m.handleRouteRecalculation(ctx, client, session)
			m.sendIncident(client, incident, action)
		} else {
			m.sendIncident(client, incident, action)
		}
	}
}

// isIncidentOnRoute returns true if an incident is on the current route.
func (m *Multicaster) isIncidentOnRoute(incident *Incident, session *navigation.Session) bool {
	return gis.IsPointInPolyline(
		gis.Point{Lat: incident.Lat, Lon: incident.Lon},
		convertNavPointsToGIS(session.Route.Polyline),
		30,
	)
}

// handleRouteRecalculation handles the route recalculation and notifies the client.
func (m *Multicaster) handleRouteRecalculation(ctx context.Context, client *ws.Client, session *navigation.Session) {
	origin := navigation.Location{
		Lat: session.LastPosition.Lat,
		Lon: session.LastPosition.Lon,
	}
	session.Route.Locations[0] = origin

	alternates := 0
	req := routing.RouteRequest{
		Locations:  convertLocationsToLocationRequests(session.Route.Locations),
		Costing:    "auto",
		Alternates: &alternates,
	}

	newRoute, err := m.RoutingClient.CalculateRoute(ctx, req)
	if err != nil {
		log.Println(err)
		return
	}

	var newPolyline []navigation.Point
	for _, leg := range newRoute.Legs {
		newPolyline = append(newPolyline, leg.Shape...)
	}

	// Met à jour la session (optionnel)
	session.Route.Polyline = newPolyline
	session.UpdatedAt = time.Now()
	if err = m.SessionCache.SetSession(ctx, session); err != nil {
		log.Printf("failed to save session to cache: %v", err)
	}

	payload, _ := json.Marshal(struct {
		Route *routing.Route `json:"route"`
		Info  string         `json:"info"`
	}{
		Route: newRoute,
		Info:  "recalculated_due_to_incident",
	})
	client.Send(ws.Message{
		Type: "route",
		Data: payload,
	})
}

// sendIncident sends a single incident to the client.
func (m *Multicaster) sendIncident(client *ws.Client, incident *Incident, action string) {
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

func convertLocationsToLocationRequests(points []navigation.Location) []routing.LocationRequest {
	res := make([]routing.LocationRequest, len(points))
	for i, p := range points {
		res[i] = routing.LocationRequest{
			Lat: p.Lat,
			Lon: p.Lon,
		}
	}
	return res
}

func convertNavPointsToGIS(points []navigation.Point) []gis.Point {
	res := make([]gis.Point, len(points))
	for i, p := range points {
		res[i] = gis.Point{Lat: p.Lat, Lon: p.Lon}
	}
	return res
}
