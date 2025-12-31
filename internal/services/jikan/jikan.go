package jikan

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.jikan.moe/v4"

type SearchResponse[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		HasNextPage bool `json:"has_next_page"`
	} `json:"pagination"`
}

type Anime struct {
	MalID   int    `json:"mal_id"`
	Title   string `json:"title"`
	Images  struct {
		JPG  struct{ ImageURL string `json:"image_url"` }  `json:"jpg"`
		Webp struct{ ImageURL string `json:"image_url"` } `json:"webp"`
	} `json:"images"`
	Synopsis string `json:"synopsis"`
	Type     string `json:"type"`
	Genres   []struct{ Name string `json:"name"` } `json:"genres"`
	Aired    struct {
		From *time.Time `json:"from"`
	} `json:"aired"`
	Score    *float32 `json:"score"`
	ScoredBy *int     `json:"scored_by"`
}

type Manga struct {
	MalID     int    `json:"mal_id"`
	Title     string `json:"title"`
	Images    struct {
		JPG  struct{ ImageURL string `json:"image_url"` }  `json:"jpg"`
		Webp struct{ ImageURL string `json:"image_url"` } `json:"webp"`
	} `json:"images"`
	Synopsis  string `json:"synopsis"`
	Genres    []struct{ Name string `json:"name"` } `json:"genres"`
	Published struct {
		From *time.Time `json:"from"`
	} `json:"published"`
	Score    *float32 `json:"score"`
	ScoredBy *int     `json:"scored_by"`
}

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) doJSON(ctx context.Context, rawURL string, out any) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	req.Header.Set("User-Agent", "completionist-api-go/1.0 (+https://localhost)")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("jikan %s: %s", resp.Status, string(b))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) SearchAnime(ctx context.Context, q string, page int) (*SearchResponse[Anime], error) {
	u, _ := url.Parse(baseURL + "/anime")
	v := url.Values{}
	v.Set("q", q)
	v.Set("page", fmt.Sprintf("%d", page))
	u.RawQuery = v.Encode()

	var out SearchResponse[Anime]
	if err := c.doJSON(ctx, u.String(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SearchManga(ctx context.Context, q string, page int) (*SearchResponse[Manga], error) {
	u, _ := url.Parse(baseURL + "/manga")
	v := url.Values{}
	v.Set("q", q)
	v.Set("page", fmt.Sprintf("%d", page))
	u.RawQuery = v.Encode()

	var out SearchResponse[Manga]
	if err := c.doJSON(ctx, u.String(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetAnime(ctx context.Context, malID int) (*Anime, error) {
	u := fmt.Sprintf("%s/anime/%d", baseURL, malID)
	var wrapper struct{ Data Anime `json:"data"` }
	if err := c.doJSON(ctx, u, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

func (c *Client) GetManga(ctx context.Context, malID int) (*Manga, error) {
	u := fmt.Sprintf("%s/manga/%d", baseURL, malID)
	var wrapper struct{ Data Manga `json:"data"` }
	if err := c.doJSON(ctx, u, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
