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
	Duration      int             `json:"duration"` // in seconds
	RodeAt        time.Time       `json:"rode_at"`
	CreatedAt     time.Time       `json:"created_at"`
	Media         json.RawMessage `json:"media"`
}

type RideModel struct {
	DB *sql.DB
}

func (m *RideModel) Create(title string, description string, trailName string, distance float64, duration int, rodeAt time.Time, media json.RawMessage) (int, error) {
	stmt := `INSERT INTO rides (title, description, trail_name, distance_miles, duration, rode_at, media) VALUES (?, ?, ?, ?, ?, ?, ?)`

	// exec works in 3 steps
	// creates a prepared statement on db with provided stmt, db parses and compiles the stmt, then stores it ready for execution
	// passes params to to db, the db then executes the prepared stmt with the params
	// then closes the prepared stmt on the db
	result, err := m.DB.Exec(stmt,
		title,
		description,
		trailName,
		distance,
		duration,
		rodeAt,
		media,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *RideModel) Get(id int) (*Ride, error) {
	r := &Ride{}
	stmt := `SELECT id, title, description, trail_name, distance_miles, duration, rode_at, created_at, media FROM rides WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&r.ID, &r.Title, &r.Description, &r.TrailName, &r.DistanceMiles, &r.Duration, &r.RodeAt, &r.CreatedAt, &r.Media)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return r, nil
}

func (m *RideModel) Latest() ([]*Ride, error) {
	stmt := `SELECT id, title, description, trail_name, distance_miles, duration, rode_at, created_at, media FROM rides ORDER BY rode_at DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := []*Ride{}
	for rows.Next() {
		var ride Ride
		err := rows.Scan(&ride.ID, &ride.Title, &ride.Description, &ride.TrailName, &ride.DistanceMiles, &ride.Duration, &ride.RodeAt, &ride.CreatedAt, &ride.Media)
		if err != nil {
			return nil, err
		}
		r = append(r, &ride)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return r, nil
}
