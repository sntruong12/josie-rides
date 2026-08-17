package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Ride represents a row in the rides table.
type Ride struct {
	ID            uint            `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	TrailName     sql.NullString  `json:"trail_name"`
	DistanceMiles float64         `json:"distance_miles"`
	Duration      time.Duration   `json:"duration"`
	RodeAt        time.Time       `json:"rode_at"`
	CreatedAt     time.Time       `json:"created_at"`
	Media         json.RawMessage `json:"media"`
}

type RideModel struct {
	DB *sql.DB
}

func (m *RideModel) Create(title string, description string, trailName string, distance float64, duration time.Duration, rodeAt time.Time, media json.RawMessage) (int, error) {
	return 0, nil
}

func (m *RideModel) Get(id int) (*Ride, error) {
	return nil, nil
}

func (m *RideModel) Latest() ([]*Ride, error) {
	return nil, nil
}
