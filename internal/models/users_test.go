package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModelExists(t *testing.T) {
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
			// Set up fresh DB and model for each sub-test
			db := newTestDB(t)
			m := UserModel{DB: db}

			exists, err := m.Exists(tt.userID)

			assert.Equal(t, tt.want, exists)
			assert.NoError(t, err)
		})
	}
}
