package gis

import (
	"math"
)

// EarthRadius in meters
const EarthRadius = 6378137

// Degrees to radians conversion
const degToRad = math.Pi / 180

// Haversine distance between two points in meters
func Haversine(a, b Point) float64 {
	dLat := (b.Lat - a.Lat) * math.Pi / 180.0
	dLon := (b.Lon - a.Lon) * math.Pi / 180.0

	lat1 := a.Lat * math.Pi / 180.0
	lat2 := b.Lat * math.Pi / 180.0

	sinDlat := math.Sin(dLat / 2)
	sinDlon := math.Sin(dLon / 2)

	aVal := sinDlat*sinDlat + sinDlon*sinDlon*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(aVal), math.Sqrt(1-aVal))
	return EarthRadius * c
}

// IsPointInPolyline returns true if given point is within tolerance distance (in metres) from the polyline.
func IsPointInPolyline(point Point, polyline []Point, tolerance float64) bool {
	if len(polyline) == 0 {
		return false
	}
	if len(polyline) == 1 {
		// Polyline mono-point 😿
		return Haversine(point, polyline[0]) <= tolerance
	}

	for i := 0; i < len(polyline)-1; i++ {
		if distanceToSegment(point, polyline[i], polyline[i+1]) <= tolerance {
			return true
		}
	}
	return false
}

// distanceToSegment calculates the minimum distance (in metres) from point P to the segment [A, B].
func distanceToSegment(P, A, B Point) float64 {
	// Convert lat/lon to radians
	lat1 := A.Lat * degToRad
	lon1 := A.Lon * degToRad
	lat2 := B.Lat * degToRad
	lon2 := B.Lon * degToRad
	latP := P.Lat * degToRad
	lonP := P.Lon * degToRad

	// Use a reference latitude for more accurate projection.
	// We ignore geodesic constraints because BLC frère.
	// We could use the cross-track distance, but I don't think we need this accuracy.
	// https://www.movable-type.co.uk/scripts/latlong.html
	latRef := (lat1 + lat2) / 2
	cosLatRef := math.Cos(latRef)

	// Project points in local Cartesian coordinates (x in east-west, y in north-south)
	xA, yA := lon1*EarthRadius*cosLatRef, lat1*EarthRadius
	xB, yB := lon2*EarthRadius*cosLatRef, lat2*EarthRadius
	xP, yP := lonP*EarthRadius*cosLatRef, latP*EarthRadius

	dx, dy := xB-xA, yB-yA

	// Degenerate segment case (A == B)
	if dx == 0 && dy == 0 {
		return math.Hypot(xP-xA, yP-yA)
	}

	// Orthogonal projection of point P onto segment AB
	t := ((xP-xA)*dx + (yP-yA)*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t)) // Clamp t dans [0,1]
	xProj := xA + t*dx
	yProj := yA + t*dy

	// Euclidean distance in metres
	return math.Hypot(xP-xProj, yP-yProj)
}
