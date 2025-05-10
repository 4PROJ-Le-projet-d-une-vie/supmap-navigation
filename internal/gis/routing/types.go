package gis

import (
	"errors"
	"fmt"
	"supmap-navigation/internal/navigation"
)

type RouteRequest struct {
	Locations      []LocationRequest `json:"locations"`
	Costing        Costing           `json:"costing"`
	CostingOptions *CostingOptions   `json:"costing_options,omitempty"`
	Language       *string           `json:"language,omitempty"`
	Alternates     *int              `json:"alternates,omitempty"`
}

func (r RouteRequest) Validate() error {
	if len(r.Locations) < 2 {
		return errors.New("at least 2 locations must be provided")
	}
	if !r.Costing.IsValid() {
		return errors.New(fmt.Sprintf("costing %q is invalid", r.Costing))
	}
	return nil
}

type LocationRequest struct {
	Lat  float64       `json:"lat"`
	Lon  float64       `json:"lon"`
	Type *LocationType `json:"type,omitempty"`
	Name *string       `json:"name,omitempty"`
}

type LocationType string

const (
	LocationTypeBreak        LocationType = "break"
	LocationTypeThrough      LocationType = "through"
	LocationTypeVia          LocationType = "via"
	LocationTypeBreakThrough LocationType = "break_through"
)

func (lt LocationType) IsValid() bool {
	switch lt {
	case LocationTypeBreak, LocationTypeThrough, LocationTypeVia, LocationTypeBreakThrough:
		return true
	default:
		return false
	}
}

type Costing string

const (
	CostingAuto         Costing = "auto"
	CostingBicycle      Costing = "bicycle"
	CostingTruck        Costing = "truck"
	CostingMotorScooter Costing = "motor_scooter"
	CostingPedestrian   Costing = "pedestrian"
)

func (c Costing) IsValid() bool {
	switch c {
	case CostingAuto, CostingBicycle, CostingTruck, CostingMotorScooter, CostingPedestrian:
		return true
	default:
		return false
	}
}

// Ratio represents a float between 0 and 1.
type Ratio float64

func (r Ratio) IsValid() bool {
	if r < 0.0 || r > 1.0 {
		return false
	}
	return true
}

type CostingOptions struct {
	UseHighways *Ratio `json:"use_highways,omitempty"`
	UseTolls    *Ratio `json:"use_tolls,omitempty"`
	UseTracks   *Ratio `json:"use_tracks,omitempty"`
}

// Response specific

type Route Trip

type RouteResponse struct {
	Data    []Route `json:"data"`
	Message string  `json:"message"`
}

type Trip struct {
	Locations []LocationResponse `json:"locations"`
	Legs      []Leg              `json:"legs"`
	Summary   Summary            `json:"summary"`
}

type Maneuver struct {
	Type                uint8    `json:"type"`
	Instruction         string   `json:"instruction"`
	StreetNames         []string `json:"street_names"`
	Time                float64  `json:"time"`
	Length              float64  `json:"length"`
	BeginShapeIndex     uint     `json:"begin_shape_index"`
	EndShapeIndex       uint     `json:"end_shape_index"`
	RoundaboutExitCount *uint8   `json:"roundabout_exit_count,omitempty"`
}

type Summary struct {
	Time   float64 `json:"time"`
	Length float64 `json:"length"`
}

type Leg struct {
	Maneuvers []Maneuver         `json:"maneuvers"`
	Summary   Summary            `json:"summary"`
	Shape     []navigation.Point `json:"shape"`
}

type LocationResponse struct {
	Lat           float64      `json:"lat"`
	Lon           float64      `json:"lon"`
	Type          LocationType `json:"type"`
	OriginalIndex int          `json:"original_index"`
	Name          *string      `json:"name,omitempty"`
}
