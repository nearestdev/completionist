package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int64      `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Role         UserRole   `db:"role" json:"role"`
	XP           int64      `db:"xp" json:"xp"`
	Level        int        `db:"level" json:"level"`
	BannedAt     *time.Time `db:"banned_at" json:"bannedAt,omitempty"`
	BanReason    *string    `db:"ban_reason" json:"banReason,omitempty"`
	ProfileHTML  *string    `db:"profile_html" json:"profileHtml,omitempty"`
	ProfileTheme *json.RawMessage `db:"profile_theme" json:"profileTheme,omitempty" swaggertype:"string"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`
}
