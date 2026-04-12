package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ModerationResult struct {
	Flagged    bool            `json:"flagged"`
	Categories json.RawMessage `json:"categories"`
	Scores     json.RawMessage `json:"category_scores"`
}

type ModerationService struct {
	apiKey  string
	enabled bool
}

func NewModerationService(apiKey string, enabled bool) *ModerationService {
	return &ModerationService{apiKey: apiKey, enabled: enabled}
}

func (s *ModerationService) IsEnabled() bool {
	return s.enabled && s.apiKey != ""
}

func (s *ModerationService) CheckContent(text string) (*ModerationResult, error) {
	if !s.IsEnabled() {
		return &ModerationResult{Flagged: false}, nil
	}

	body, _ := json.Marshal(map[string]string{"input": text})
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/moderations", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("moderation API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("moderation API returned %d: %s", resp.StatusCode, string(b))
	}

	var apiResp struct {
		Results []ModerationResult `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode moderation response: %w", err)
	}
	if len(apiResp.Results) == 0 {
		return &ModerationResult{Flagged: false}, nil
	}
	return &apiResp.Results[0], nil
}
