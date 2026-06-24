package googlebooks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseAPI = "https://www.googleapis.com/books/v1"

type Client struct {
	apiKey string
	http   *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

type SearchResult struct {
	Kind       string `json:"kind"`
	TotalItems int    `json:"totalItems"`
	Items      []Book `json:"items"`
}

type Book struct {
	ID         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type VolumeInfo struct {
	Title         string     `json:"title"`
	Subtitle      string     `json:"subtitle"`
	Authors       []string   `json:"authors"`
	Publisher     string     `json:"publisher"`
	PublishedDate string     `json:"publishedDate"`
	Description   string     `json:"description"`
	PageCount     int        `json:"pageCount"`
	Categories    []string   `json:"categories"`
	ImageLinks    ImageLinks `json:"imageLinks"`
}

type ImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

func (c *Client) Search(ctx context.Context, query string) (*SearchResult, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/volumes", baseAPI))
	q := u.Query()
	q.Set("q", query)
	q.Set("key", c.apiKey)
	u.RawQuery = q.Encode()
	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google books api error: %s - %s", resp.Status, string(bodyBytes))
	}
	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetVolume(ctx context.Context, volumeID string) (*Book, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/volumes/%s", baseAPI, volumeID))
	q := u.Query()
	q.Set("key", c.apiKey)
	u.RawQuery = q.Encode()
	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google books api error: %s - %s", resp.Status, string(bodyBytes))
	}
	var book Book
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, err
	}
	return &book, nil
}
