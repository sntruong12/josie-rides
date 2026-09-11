package mocks

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sntruong12/josie-rides/internal/models"
)

var mockRide = &models.Ride{
	ID:            1,
	Title:         "Morning Ride",
	Description:   "A fun ride through the scenic trail.",
	TrailName:     sql.NullString{String: "River Trail", Valid: true},
	DistanceMiles: 12.5,
	Duration:      3600,
	RodeAt:        time.Now(),
	CreatedAt:     time.Now(),
	Media:         json.RawMessage(`[]`),
}

type RideModel struct{}

func (m *RideModel) Create(title string, description string, trailName string, distance float64, duration int, rodeAt time.Time, media json.RawMessage) (int, error) {
	return 2, nil
}

func (m *RideModel) Get(id int) (*models.Ride, error) {
	switch id {
	case 1:
		return mockRide, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *RideModel) Latest() ([]*models.Ride, error) {
	return []*models.Ride{mockRide}, nil
}
