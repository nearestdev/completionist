package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/yohcop/openid-go"
)

type Client struct {
	apiKey     string
	base       string
	http       *http.Client
	publicBase string
}

func New(apiKey, publicBase string) *Client {
	return &Client{
		apiKey:     apiKey,
		base:       "https://api.steampowered.com",
		http:       &http.Client{Timeout: 12 * time.Second},
		publicBase: publicBase,
	}
}

type OwnedGamesResponse struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			AppID                    int    `json:"appid"`
			Name                     string `json:"name"`
			PlaytimeForever          int    `json:"playtime_forever"`
			ImgIconURL               string `json:"img_icon_url"`
			ImgLogoURL               string `json:"img_logo_url"`
			HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
			PlaytimeWindowsForever   int    `json:"playtime_windows_forever"`
			PlaytimeMacForever       int    `json:"playtime_mac_forever"`
			PlaytimeLinuxForever     int    `json:"playtime_linux_forever"`
			RtimeLastPlayed          int64  `json:"rtime_last_played"`
		} `json:"games"`
	} `json:"response"`
}

type PlayerSummary struct {
	PersonaName string `json:"personaname"`
	AvatarFull  string `json:"avatarfull"`
}

type PlayerSummariesResponse struct {
	Response struct {
		Players []PlayerSummary `json:"players"`
	} `json:"response"`
}

type PlayerAchievementsResponse struct {
	Playerstats struct {
		SteamID      string `json:"steamID"`
		GameName     string `json:"gameName"`
		Achievements []struct {
			APIName    string `json:"apiname"`
			Achieved   int    `json:"achieved"`
			UnlockTime int64  `json:"unlocktime"`
		} `json:"achievements"`
		Success bool `json:"success"`
	} `json:"playerstats"`
}

type SchemaForGameResponse struct {
	Game struct {
		GameName           string `json:"gameName"`
		GameVersion        string `json:"gameVersion"`
		AvailableGameStats struct {
			Achievements []struct {
				Name         string `json:"name"`
				DefaultValue int    `json:"defaultvalue"`
				DisplayName  string `json:"displayName"`
				Hidden       int    `json:"hidden"`
				Icon         string `json:"icon"`
				IconGray     string `json:"icongray"`
			} `json:"achievements"`
		} `json:"availableGameStats"`
	} `json:"game"`
}

func (c *Client) GetOwnedGames(ctx context.Context, steamID string, includeAppInfo, includeFree bool) (*OwnedGamesResponse, error) {
	u := fmt.Sprintf("%s/IPlayerService/GetOwnedGames/v1", c.base)
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("steamid", steamID)
	if includeAppInfo {
		q.Set("include_appinfo", "1")
	}
	if includeFree {
		q.Set("include_played_free_games", "1")
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", u+"?"+q.Encode(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out OwnedGamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetPlayerSummary(ctx context.Context, steamID string) (*PlayerSummary, error) {
	u := fmt.Sprintf("%s/ISteamUser/GetPlayerSummaries/v2", c.base)
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("steamids", steamID)
	req, _ := http.NewRequestWithContext(ctx, "GET", u+"?"+q.Encode(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out PlayerSummariesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Response.Players) == 0 {
		return nil, fmt.Errorf("no player")
	}
	return &out.Response.Players[0], nil
}

func (c *Client) GetPlayerAchievements(ctx context.Context, steamID string, appID int) (*PlayerAchievementsResponse, error) {
	u := fmt.Sprintf("%s/ISteamUserStats/GetPlayerAchievements/v1", c.base)
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("steamid", steamID)
	q.Set("appid", strconv.Itoa(appID))
	req, _ := http.NewRequestWithContext(ctx, "GET", u+"?"+q.Encode(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out PlayerAchievementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetSchemaForGame(ctx context.Context, appID int) (*SchemaForGameResponse, error) {
	u := fmt.Sprintf("%s/ISteamUserStats/GetSchemaForGame/v2", c.base)
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("appid", strconv.Itoa(appID))
	req, _ := http.NewRequestWithContext(ctx, "GET", u+"?"+q.Encode(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out SchemaForGameResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) BuildOpenIDRedirect(authToken string) string {
	u := url.URL{
		Scheme: "https",
		Host:   "steamcommunity.com",
		Path:   "/openid/login",
	}
	q := url.Values{}

	returnTo := c.publicBase + "/api/auth/steam/callback"
	if authToken != "" {
		returnTo = returnTo + "?auth_token=" + url.QueryEscape(authToken)
	}

	q.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	q.Set("openid.mode", "checkid_setup")
	q.Set("openid.return_to", returnTo)
	q.Set("openid.realm", c.publicBase)
	q.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	q.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) VerifyOpenID(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	fullURL := c.publicBase + r.URL.Path
	if r.URL.RawQuery != "" {
		fullURL += "?" + r.URL.RawQuery
	}
	claimedID, err := openid.Verify(fullURL, openid.NewSimpleDiscoveryCache(), openid.NewSimpleNonceStore())
	if err != nil {
		return "", fmt.Errorf("openid verify failed: %w", err)
	}
	if claimedID == "" {
		return "", fmt.Errorf("empty claimed id")
	}
	last := claimedID[strings.LastIndex(claimedID, "/")+1:]
	if last == "" {
		return "", fmt.Errorf("invalid claimed id")
	}
	return last, nil
}