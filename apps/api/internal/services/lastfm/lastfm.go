package lastfm

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL     = "https://ws.audioscrobbler.com/2.0/"
	authURLBase = "https://www.last.fm/api/auth/"
)

type Client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
	baseURL    string
}

func New(apiKey, apiSecret, baseURL string) *Client {
	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}

func (c *Client) createSignature(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sigStr strings.Builder
	for _, k := range keys {
		sigStr.WriteString(k)
		sigStr.WriteString(params[k])
	}
	sigStr.WriteString(c.apiSecret)
	hasher := md5.New()
	hasher.Write([]byte(sigStr.String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *Client) doRequest(ctx context.Context, params map[string]string, result interface{}) error {
	params["api_key"] = c.apiKey
	params["format"] = "json"

	isSigned := params["api_sig"] != ""
	if isSigned {
		delete(params, "format")
		sig := c.createSignature(params)
		params["api_sig"] = sig
		params["format"] = "json"
	}

	u, _ := url.Parse(baseURL)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("last.fm API error: %s - %s", resp.Status, string(bodyBytes))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) GetAuthURL(authToken string) string {
	return fmt.Sprintf("%s?api_key=%s&cb=%s/api/auth/lastfm/callback?auth_token=%s",
		authURLBase, c.apiKey, c.baseURL, url.QueryEscape(authToken))
}

type SessionResponse struct {
	Session struct {
		Name       string `json:"name"`
		Key        string `json:"key"`
		Subscriber int    `json:"subscriber"`
	} `json:"session"`
}

func (c *Client) GetSession(ctx context.Context, token string) (*SessionResponse, error) {
	params := map[string]string{
		"method":  "auth.getSession",
		"token":   token,
		"api_sig": "true",
	}
	var sessionResp SessionResponse
	if err := c.doRequest(ctx, params, &sessionResp); err != nil {
		return nil, err
	}
	if sessionResp.Session.Name == "" {
		return nil, fmt.Errorf("failed to get session from last.fm")
	}
	return &sessionResp, nil
}

type RecentTracksResponse struct {
	RecentTracks struct {
		Track []Track `json:"track"`
	} `json:"recenttracks"`
}

type Track struct {
	Artist struct {
		Name string `json:"#text"`
	} `json:"artist"`
	Name  string `json:"name"`
	Album struct {
		Name string `json:"#text"`
	} `json:"album"`
	Image []struct {
		URL  string `json:"#text"`
		Size string `json:"size"`
	} `json:"image"`
	URL        string `json:"url"`
	NowPlaying *struct {
		NowPlaying string `json:"nowplaying"`
	} `json:"@attr,omitempty"`
}

func (c *Client) GetRecentTracks(ctx context.Context, user string, limit int) (*RecentTracksResponse, error) {
	params := map[string]string{
		"method": "user.getrecenttracks",
		"user":   user,
		"limit":  strconv.Itoa(limit),
	}
	var recentTracksResp RecentTracksResponse
	if err := c.doRequest(ctx, params, &recentTracksResp); err != nil {
		return nil, err
	}
	return &recentTracksResp, nil
}
