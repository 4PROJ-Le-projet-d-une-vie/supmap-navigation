package incidents

import (
	"time"
)

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
