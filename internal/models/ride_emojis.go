package models

import (
	"database/sql"
	"time"
)

type RideEmojiModelInterface interface {
	CreateOrUpdate(rideID, userID int, emoji string) error
	Delete(id, rideID, userID int) error
}

type RideEmoji struct {
	ID        int
	RideID    int
	UserID    int
	Emoji     string
	CreatedAt time.Time
}

type RideEmojiModel struct {
	DB *sql.DB
}

// We'll use the CreateOrUpdate method to add a new record to the "users" table.
func (m *RideEmojiModel) CreateOrUpdate(rideID, userID int, emoji string) error {
	stmt := `INSERT INTO ride_emojis (ride_id, user_id, emoji, created_at)
		VALUES(?, ?, ?, UTC_TIMESTAMP()) as new_emoji
		ON DUPLICATE KEY UPDATE
			emoji = VALUES(new_emoji.emoji);`
	_, err := m.DB.Exec(stmt, rideID, userID, emoji)
	if err != nil {
		return err
	}
	return nil
}

func (m *RideEmojiModel) Delete(id, rideID, userID int) error {
	stmt := `DELETE FROM ride_emojis WHERE id = ? AND ride_id = ? AND user_id = ?;`
	_, err := m.DB.Exec(stmt, id, rideID, userID)
	if err != nil {
		return err
	}
	return nil
}
