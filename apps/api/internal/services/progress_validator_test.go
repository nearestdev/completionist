package services

import (
	"encoding/json"
	"testing"

	"github.com/nearestdev/completionist/internal/models"
)

func ptrInt(v int) *int           { return &v }
func ptrFloat(v float64) *float64 { return &v }
func ptrStr(v string) *string     { return &v }
func ptrBool(v bool) *bool        { return &v }

func raw(s string) *json.RawMessage {
	m := json.RawMessage(s)
	return &m
}

func rawFrom(t *testing.T, data models.ProgressData) *json.RawMessage {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal progress data: %v", err)
	}
	m := json.RawMessage(b)
	return &m
}

func TestValidateEmptyProgress(t *testing.T) {
	pv := NewProgressValidator()

	tests := []struct {
		name string
		data *json.RawMessage
	}{
		{"nil pointer", nil},
		{"empty bytes", raw("")},
		{"empty object", raw("{}")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, it := range []models.ItemType{
				models.ItemTypeBook,
				models.ItemTypeManga,
				models.ItemTypeGame,
				models.ItemTypeSeries,
				models.ItemTypeMovie,
				models.ItemTypeMusic,
			} {
				if err := pv.Validate(tt.data, it); err != nil {
					t.Errorf("Validate(%s) = %v, want nil", it, err)
				}
			}
		})
	}
}

func TestValidateMalformedJSON(t *testing.T) {
	pv := NewProgressValidator()

	if err := pv.Validate(raw("{not json"), models.ItemTypeBook); err == nil {
		t.Fatal("expected error for malformed json, got nil")
	}
}

func TestValidateUnsupportedItemType(t *testing.T) {
	pv := NewProgressValidator()

	err := pv.Validate(raw(`{"current":1,"total":2}`), models.ItemType("podcast"))
	if err == nil {
		t.Fatal("expected error for unsupported item type, got nil")
	}
}

func TestValidatePagedTypes(t *testing.T) {
	pv := NewProgressValidator()

	tests := []struct {
		name     string
		itemType models.ItemType
		data     models.ProgressData
		wantErr  bool
	}{
		{"book valid pages", models.ItemTypeBook, models.ProgressData{Current: ptrInt(50), Total: ptrInt(100), Unit: ptrStr("pages")}, false},
		{"book valid chapters", models.ItemTypeBook, models.ProgressData{Current: ptrInt(3), Total: ptrInt(10), Unit: ptrStr("chapters")}, false},
		{"book no unit", models.ItemTypeBook, models.ProgressData{Current: ptrInt(3), Total: ptrInt(10)}, false},
		{"book invalid unit", models.ItemTypeBook, models.ProgressData{Current: ptrInt(3), Total: ptrInt(10), Unit: ptrStr("volumes")}, true},
		{"book negative current", models.ItemTypeBook, models.ProgressData{Current: ptrInt(-1), Total: ptrInt(10)}, true},
		{"book total zero", models.ItemTypeBook, models.ProgressData{Current: ptrInt(0), Total: ptrInt(0)}, true},
		{"book total negative", models.ItemTypeBook, models.ProgressData{Current: ptrInt(1), Total: ptrInt(-5)}, true},
		{"book current exceeds total", models.ItemTypeBook, models.ProgressData{Current: ptrInt(11), Total: ptrInt(10)}, true},
		{"book current equals total", models.ItemTypeBook, models.ProgressData{Current: ptrInt(10), Total: ptrInt(10)}, false},
		{"book only current set skips checks", models.ItemTypeBook, models.ProgressData{Current: ptrInt(-99)}, false},
		{"book only total set skips checks", models.ItemTypeBook, models.ProgressData{Total: ptrInt(-99)}, false},

		{"manga valid volumes", models.ItemTypeManga, models.ProgressData{Current: ptrInt(2), Total: ptrInt(12), Unit: ptrStr("volumes")}, false},
		{"manga valid chapters", models.ItemTypeManga, models.ProgressData{Current: ptrInt(40), Total: ptrInt(120), Unit: ptrStr("chapters")}, false},
		{"manga invalid unit pages", models.ItemTypeManga, models.ProgressData{Current: ptrInt(1), Total: ptrInt(10), Unit: ptrStr("pages")}, true},
		{"manga negative current", models.ItemTypeManga, models.ProgressData{Current: ptrInt(-3), Total: ptrInt(10)}, true},
		{"manga total zero", models.ItemTypeManga, models.ProgressData{Current: ptrInt(1), Total: ptrInt(0)}, true},
		{"manga current exceeds total", models.ItemTypeManga, models.ProgressData{Current: ptrInt(13), Total: ptrInt(12)}, true},

		{"series valid episodes", models.ItemTypeSeries, models.ProgressData{Current: ptrInt(5), Total: ptrInt(24), Unit: ptrStr("episodes")}, false},
		{"series no unit", models.ItemTypeSeries, models.ProgressData{Current: ptrInt(5), Total: ptrInt(24)}, false},
		{"series invalid unit", models.ItemTypeSeries, models.ProgressData{Current: ptrInt(5), Total: ptrInt(24), Unit: ptrStr("chapters")}, true},
		{"series negative current", models.ItemTypeSeries, models.ProgressData{Current: ptrInt(-1), Total: ptrInt(24)}, true},
		{"series total zero", models.ItemTypeSeries, models.ProgressData{Current: ptrInt(0), Total: ptrInt(0)}, true},
		{"series current exceeds total", models.ItemTypeSeries, models.ProgressData{Current: ptrInt(25), Total: ptrInt(24)}, true},
		{"series season valid", models.ItemTypeSeries, models.ProgressData{Season: ptrInt(2)}, false},
		{"series season one", models.ItemTypeSeries, models.ProgressData{Season: ptrInt(1)}, false},
		{"series season zero", models.ItemTypeSeries, models.ProgressData{Season: ptrInt(0)}, true},
		{"series season negative", models.ItemTypeSeries, models.ProgressData{Season: ptrInt(-1)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pv.Validate(rawFrom(t, tt.data), tt.itemType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGame(t *testing.T) {
	pv := NewProgressValidator()

	tests := []struct {
		name    string
		data    models.ProgressData
		wantErr bool
	}{
		{"hours valid", models.ProgressData{HoursPlayed: ptrFloat(12.5)}, false},
		{"hours zero", models.ProgressData{HoursPlayed: ptrFloat(0)}, false},
		{"hours negative", models.ProgressData{HoursPlayed: ptrFloat(-0.1)}, true},
		{"achievements valid", models.ProgressData{AchievementsUnlocked: ptrInt(10), AchievementsTotal: ptrInt(50)}, false},
		{"achievements equal total", models.ProgressData{AchievementsUnlocked: ptrInt(50), AchievementsTotal: ptrInt(50)}, false},
		{"achievements negative", models.ProgressData{AchievementsUnlocked: ptrInt(-1), AchievementsTotal: ptrInt(50)}, true},
		{"achievements total zero", models.ProgressData{AchievementsUnlocked: ptrInt(0), AchievementsTotal: ptrInt(0)}, true},
		{"achievements total negative", models.ProgressData{AchievementsUnlocked: ptrInt(1), AchievementsTotal: ptrInt(-5)}, true},
		{"achievements exceed total", models.ProgressData{AchievementsUnlocked: ptrInt(51), AchievementsTotal: ptrInt(50)}, true},
		{"achievements unlocked only no total", models.ProgressData{AchievementsUnlocked: ptrInt(9999)}, false},
		{"achievements total only is ignored", models.ProgressData{AchievementsTotal: ptrInt(-5)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pv.Validate(rawFrom(t, tt.data), models.ItemTypeGame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMovieIsNoOp(t *testing.T) {
	pv := NewProgressValidator()

	cases := []*json.RawMessage{
		raw(`{"watched":true}`),
		raw(`{"watched":false}`),
		raw(`{"current":-50,"total":-10,"unit":"nonsense"}`),
		raw(`{"season":-99}`),
	}

	for _, c := range cases {
		if err := pv.Validate(c, models.ItemTypeMovie); err != nil {
			t.Errorf("Validate(movie, %s) = %v, want nil", string(*c), err)
		}
	}
}

func TestValidateMusic(t *testing.T) {
	pv := NewProgressValidator()

	tests := []struct {
		name    string
		data    models.ProgressData
		wantErr bool
	}{
		{"tracks valid", models.ProgressData{TracksListened: ptrInt(5), TotalTracks: ptrInt(12)}, false},
		{"tracks equal total", models.ProgressData{TracksListened: ptrInt(12), TotalTracks: ptrInt(12)}, false},
		{"tracks negative", models.ProgressData{TracksListened: ptrInt(-1), TotalTracks: ptrInt(12)}, true},
		{"tracks total zero", models.ProgressData{TracksListened: ptrInt(0), TotalTracks: ptrInt(0)}, true},
		{"tracks total negative", models.ProgressData{TracksListened: ptrInt(1), TotalTracks: ptrInt(-3)}, true},
		{"tracks exceed total", models.ProgressData{TracksListened: ptrInt(13), TotalTracks: ptrInt(12)}, true},
		{"tracks listened only no total", models.ProgressData{TracksListened: ptrInt(9999)}, false},
		{"albums valid", models.ProgressData{AlbumsListened: ptrInt(2), TotalAlbums: ptrInt(5)}, false},
		{"albums negative", models.ProgressData{AlbumsListened: ptrInt(-1), TotalAlbums: ptrInt(5)}, true},
		{"albums total zero", models.ProgressData{AlbumsListened: ptrInt(0), TotalAlbums: ptrInt(0)}, true},
		{"albums exceed total", models.ProgressData{AlbumsListened: ptrInt(6), TotalAlbums: ptrInt(5)}, true},
		{"albums listened only no total", models.ProgressData{AlbumsListened: ptrInt(9999)}, false},
		{"tracks and albums both valid", models.ProgressData{TracksListened: ptrInt(3), TotalTracks: ptrInt(10), AlbumsListened: ptrInt(1), TotalAlbums: ptrInt(2)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pv.Validate(rawFrom(t, tt.data), models.ItemTypeMusic)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculatePercentageEmpty(t *testing.T) {
	pv := NewProgressValidator()

	tests := []struct {
		name string
		data *json.RawMessage
	}{
		{"nil pointer", nil},
		{"empty bytes", raw("")},
		{"empty object", raw("{}")},
		{"malformed json", raw("{bad")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pv.CalculatePercentage(tt.data, models.ItemTypeBook); got != 0.0 {
				t.Errorf("CalculatePercentage() = %v, want 0", got)
			}
		})
	}
}

func TestCalculatePercentage(t *testing.T) {
	pv := NewProgressValidator()

	tests := []struct {
		name     string
		itemType models.ItemType
		data     models.ProgressData
		want     float64
	}{
		{"book half", models.ItemTypeBook, models.ProgressData{Current: ptrInt(5), Total: ptrInt(10)}, 50.0},
		{"book full", models.ItemTypeBook, models.ProgressData{Current: ptrInt(10), Total: ptrInt(10)}, 100.0},
		{"book quarter", models.ItemTypeBook, models.ProgressData{Current: ptrInt(1), Total: ptrInt(4)}, 25.0},
		{"book zero progress", models.ItemTypeBook, models.ProgressData{Current: ptrInt(0), Total: ptrInt(10)}, 0.0},
		{"book missing total", models.ItemTypeBook, models.ProgressData{Current: ptrInt(5)}, 0.0},
		{"book total zero", models.ItemTypeBook, models.ProgressData{Current: ptrInt(5), Total: ptrInt(0)}, 0.0},
		{"manga two thirds", models.ItemTypeManga, models.ProgressData{Current: ptrInt(2), Total: ptrInt(8)}, 25.0},
		{"series fifth", models.ItemTypeSeries, models.ProgressData{Current: ptrInt(3), Total: ptrInt(15)}, 20.0},

		{"game achievements", models.ItemTypeGame, models.ProgressData{AchievementsUnlocked: ptrInt(25), AchievementsTotal: ptrInt(50)}, 50.0},
		{"game achievements full", models.ItemTypeGame, models.ProgressData{AchievementsUnlocked: ptrInt(50), AchievementsTotal: ptrInt(50)}, 100.0},
		{"game missing total", models.ItemTypeGame, models.ProgressData{AchievementsUnlocked: ptrInt(25)}, 0.0},
		{"game hours only", models.ItemTypeGame, models.ProgressData{HoursPlayed: ptrFloat(40)}, 0.0},

		{"movie watched", models.ItemTypeMovie, models.ProgressData{Watched: ptrBool(true)}, 100.0},
		{"movie not watched", models.ItemTypeMovie, models.ProgressData{Watched: ptrBool(false)}, 0.0},

		{"music tracks", models.ItemTypeMusic, models.ProgressData{TracksListened: ptrInt(3), TotalTracks: ptrInt(12)}, 25.0},
		{"music albums", models.ItemTypeMusic, models.ProgressData{AlbumsListened: ptrInt(1), TotalAlbums: ptrInt(4)}, 25.0},
		{"music tracks take priority over albums", models.ItemTypeMusic, models.ProgressData{TracksListened: ptrInt(6), TotalTracks: ptrInt(10), AlbumsListened: ptrInt(1), TotalAlbums: ptrInt(4)}, 60.0},
		{"music albums when tracks total missing", models.ItemTypeMusic, models.ProgressData{TracksListened: ptrInt(6), AlbumsListened: ptrInt(1), TotalAlbums: ptrInt(2)}, 50.0},
		{"music no data", models.ItemTypeMusic, models.ProgressData{HoursPlayed: ptrFloat(1)}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pv.CalculatePercentage(rawFrom(t, tt.data), tt.itemType)
			if got != tt.want {
				t.Errorf("CalculatePercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculatePercentageUnsupportedType(t *testing.T) {
	pv := NewProgressValidator()

	got := pv.CalculatePercentage(raw(`{"current":5,"total":10}`), models.ItemType("podcast"))
	if got != 0.0 {
		t.Errorf("CalculatePercentage() = %v, want 0", got)
	}
}
