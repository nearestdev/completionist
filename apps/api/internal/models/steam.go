package models

type SteamAccount struct {
	UserID  int64  `db:"user_id" json:"userId"`
	SteamID string `db:"steam_id" json:"steamId"`
	Persona string `db:"persona" json:"persona"`
	Avatar  string `db:"avatar" json:"avatar"`
}
type AttachSteamRequest struct {
	SteamID string `json:"steamId"`
}
