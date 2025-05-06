package navigation

import (
	"context"
	"time"
)

type Session struct {
	UserID       string    `json:"user_id"`
	LastPosition Position  `json:"last_position"`
	Route        Route     `json:"route"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Position struct {
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Timestamp time.Time `json:"timestamp"`
}

type Point struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type Route struct {
	Polyline  []Point `json:"polyline"`
	Locations []Point `json:"locations"`
}

type SessionCache interface {
	SetSession(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, userID string) (*Session, error)
	DeleteSession(ctx context.Context, userID string) error
}
