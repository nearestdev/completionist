package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const baseAPI = "https://api.themoviedb.org/3"

type Client struct {
	apiKey      string
	http        *http.Client
	cfgOnce     sync.Once
	cfgErr      error
	imageBase   string
	posterSizes []string
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

type tmdbConfiguration struct {
	Images struct {
		BaseURL       string   `json:"base_url"`
		SecureBaseURL string   `json:"secure_base_url"`
		PosterSizes   []string `json:"poster_sizes"`
	} `json:"images"`
}

func (c *Client) ensureConfig(ctx context.Context) error {
	c.cfgOnce.Do(func() {
		req, _ := http.NewRequestWithContext(ctx, "GET", baseAPI+"/configuration", nil)
		q := req.URL.Query()
		q.Set("api_key", c.apiKey)
		req.URL.RawQuery = q.Encode()
		resp, err := c.http.Do(req)
		if err != nil {
			c.cfgErr = err
			return
		}
		defer resp.Body.Close()
		var cfg tmdbConfiguration
		if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
			c.cfgErr = err
			return
		}
		base := cfg.Images.SecureBaseURL
		if base == "" {
			base = cfg.Images.BaseURL
		}
		c.imageBase = base
		c.posterSizes = cfg.Images.PosterSizes
	})
	return c.cfgErr
}

func (c *Client) fullPosterURL(ctx context.Context, path string) *string {
	if path == "" {
		return nil
	}
	if err := c.ensureConfig(ctx); err != nil {
		return nil
	}
	size := "w500"
	if len(c.posterSizes) > 0 {
		size = c.posterSizes[0]
		for _, s := range c.posterSizes {
			if s == "w500" {
				size = s
				break
			}
		}
	}
	u := fmt.Sprintf("%s%s/%s", c.imageBase, size, path)
	return &u
}

type SearchResult struct {
	Page         int          `json:"page"`
	Results      []SearchItem `json:"results"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

type SearchItem struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	MediaType    string `json:"media_type"`
	GenreIDs     []int  `json:"genre_ids"`
	ReleaseDate  string `json:"release_date"`
	FirstAirDate string `json:"first_air_date"`
}

func (c *Client) SearchMovies(ctx context.Context, query string, page int) (*SearchResult, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", baseAPI+"/search/movie", nil)
	q := req.URL.Query()
	q.Set("api_key", c.apiKey)
	q.Set("query", query)
	q.Set("page", strconv.Itoa(page))
	req.URL.RawQuery = q.Encode()
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

func (c *Client) SearchTV(ctx context.Context, query string, page int) (*SearchResult, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", baseAPI+"/search/tv", nil)
	q := req.URL.Query()
	q.Set("api_key", c.apiKey)
	q.Set("query", query)
	q.Set("page", strconv.Itoa(page))
	req.URL.RawQuery = q.Encode()
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

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Genres      []Genre `json:"genres"`
}

type TV struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	FirstAirDate string  `json:"first_air_date"`
	Genres       []Genre `json:"genres"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetMovie(ctx context.Context, id int) (*Movie, *string, error) {
	u := fmt.Sprintf("%s/movie/%d", baseAPI, id)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	q := req.URL.Query()
	q.Set("api_key", c.apiKey)
	req.URL.RawQuery = q.Encode()
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	var mv Movie
	if err := json.NewDecoder(resp.Body).Decode(&mv); err != nil {
		return nil, nil, err
	}
	poster := c.fullPosterURL(ctx, mv.PosterPath)
	return &mv, poster, nil
}

func (c *Client) GetTV(ctx context.Context, id int) (*TV, *string, error) {
	u := fmt.Sprintf("%s/tv/%d", baseAPI, id)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	q := req.URL.Query()
	q.Set("api_key", c.apiKey)
	req.URL.RawQuery = q.Encode()
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	var tv TV
	if err := json.NewDecoder(resp.Body).Decode(&tv); err != nil {
		return nil, nil, err
	}
	poster := c.fullPosterURL(ctx, tv.PosterPath)
	return &tv, poster, nil
}

func (c *Client) BuildImageURL(path, size string) *string {
	if path == "" {
		return nil
	}
	if size == "" {
		size = "w500"
	}
	u := url.URL{Scheme: "https", Host: "image.tmdb.org", Path: "/t/p/" + size + path}
	s := u.String()
	return &s
}
