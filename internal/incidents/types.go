package incidents

import (
	"fmt"
)

type Incident struct {
	ID     int     `json:"id"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	TypeID int     `json:"type_id"`
}

func (i *Incident) Validate() error {
	if i.ID <= 0 {
		return fmt.Errorf("invalid ID: %d", i.ID)
	}
	if i.TypeID <= 0 {
		return fmt.Errorf("invalid TypeID: %d", i.TypeID)
	}
	if i.Lat < -90 || i.Lat > 90 {
		return fmt.Errorf("invalid latitude: %f", i.Lat)
	}
	if i.Lon < -180 || i.Lon > 180 {
		return fmt.Errorf("invalid longitude: %f", i.Lon)
	}
	return nil
}
