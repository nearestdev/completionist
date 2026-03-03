package models

type LastFMAccount struct {
	UserID     int64  `db:"user_id" json:"userId"`
	Username   string `db:"username" json:"username"`
	SessionKey string `db:"session_key" json:"sessionKey"`
	Subscriber int    `db:"subscriber" json:"subscriber"`
}