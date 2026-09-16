package models

import (
	"testing"
	"time"

	"github.com/sntruong12/josie-rides/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	// Skip the test if the "-short" flag is provided when running the test.
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to
			// our test database. Calling this here -- inside t.Run() -- means
			// that fresh database tables and data will be set up and torn down
			// for each sub-test.
			db := newTestDB(t)

			// Create a new instance of the UserModel.
			m := UserModel{db}

			// Call the UserModel.Exists() method and check that the return
			// value and error match the expected values for the sub-test.
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *User
		expectedErr error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &User{
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectedErr: nil,
		},
		{
			name:        "Zero ID",
			userID:      0,
			wantUser:    nil,
			expectedErr: ErrNoRecord,
		},
		{
			name:        "Non-existent ID",
			userID:      2,
			wantUser:    nil,
			expectedErr: ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			user, err := m.Get(tt.userID)

			assert.Equal(t, err, tt.expectedErr)

			if tt.wantUser != nil {
				assert.Equal(t, user.Name, tt.wantUser.Name)
				assert.Equal(t, user.Email, tt.wantUser.Email)
				assert.Equal(t, user.Created, tt.wantUser.Created)
			}
		})
	}
}

