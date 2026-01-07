package rawg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)
type Client struct {
	apiKey string
	http   *http.Client
	base   string
}
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 12 * time.Second},
		base:   "https://api.rawg.io/api",
	}
}
type SearchResult struct {
	Count int `json:"count"`
	Results []struct {
		ID int `json:"id"`
		Name string `json:"name"`
		Released string `json:"released"`
		BackgroundImage string `json:"background_image"`
		Rating float64 `json:"rating"`
		Platforms []struct {
			Platform struct{ ID int `json:"id"`; Name string `json:"name"` } `json:"platform"`
		} `json:"platforms"`
	} `json:"results"`
}

type Game struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Released        string `json:"released"`
	BackgroundImage string `json:"background_image"`
	Description     string `json:"description_raw"`
	Rating          float64 `json:"rating"`
	Genres          []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

type AchievementsResult struct {
	Count int `json:"count"`
	Results []struct {
		ID int `json:"id"`
		Name string `json:"name"`
		Description string `json:"description"`
		Image string `json:"image"`
		Percent float64 `json:"percent"`
	} `json:"results"`
}
func (c *Client) SearchGames(ctx context.Context, q string, page int) (*SearchResult, error) {
	u := fmt.Sprintf("%s/games", c.base)
	v := url.Values{}
	v.Set("key", c.apiKey)
	v.Set("search", q)
	v.Set("page", strconv.Itoa(page))
	req, _ := http.NewRequestWithContext(ctx, "GET", u+"?"+v.Encode(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetGame(ctx context.Context, id int) (*Game, error) {
	u := fmt.Sprintf("%s/games/%d", c.base, id)
	v := url.Values{}
	v.Set("key", c.apiKey)
	req, _ := http.NewRequestWithContext(ctx, "GET", u+"?"+v.Encode(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var game Game
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return nil, err
	}
	return &game, nil
}

func (c *Client) GameAchievements(ctx context.Context, id int, page int) (*AchievementsResult, error) {
	u := fmt.Sprintf("%s/games/%d/achievements", c.base, id)
	v := url.Values{}
	v.Set("key", c.apiKey)
	v.Set("page", strconv.Itoa(page))
	req, _ := http.NewRequestWithContext(ctx, "GET", u+"?"+v.Encode(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out AchievementsResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
