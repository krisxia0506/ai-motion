package sora

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrAPIKeyRequired  = errors.New("sora api key is required")
	ErrBaseURLRequired = errors.New("sora base url is required")
	ErrInvalidResponse = errors.New("invalid response from sora api")
	ErrRateLimited     = errors.New("rate limited by sora api")
	ErrVideoProcessing = errors.New("video is still processing")
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, ErrAPIKeyRequired
	}
	if baseURL == "" {
		return nil, ErrBaseURLRequired
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

type ImageToVideoRequest struct {
	ImageURL string
	Prompt   string
	Duration int
	Width    int
	Height   int
}

type TextToVideoRequest struct {
	Prompt   string
	Duration int
	Width    int
	Height   int
}

type VideoGenerationResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

func (c *Client) ImageToVideo(ctx context.Context, req ImageToVideoRequest) (string, error) {
	payload := map[string]interface{}{
		"prompt": req.Prompt,
	}

	if req.ImageURL != "" {
		payload["image"] = req.ImageURL
	}

	if req.Duration > 0 {
		payload["duration"] = req.Duration
	}

	if req.Width > 0 && req.Height > 0 {
		payload["size"] = fmt.Sprintf("%dx%d", req.Width, req.Height)
	}

	result, err := c.makeRequest(ctx, "videos", payload)
	if err != nil {
		return "", err
	}

	return c.extractVideoID(result)
}

func (c *Client) TextToVideo(ctx context.Context, req TextToVideoRequest) (string, error) {
	payload := map[string]interface{}{
		"prompt": req.Prompt,
	}

	if req.Duration > 0 {
		payload["duration"] = req.Duration
	}

	if req.Width > 0 && req.Height > 0 {
		payload["size"] = fmt.Sprintf("%dx%d", req.Width, req.Height)
	}

	result, err := c.makeRequest(ctx, "videos", payload)
	if err != nil {
		return "", err
	}

	return c.extractVideoID(result)
}

func (c *Client) GetVideoStatus(ctx context.Context, videoID string) (*VideoGenerationResponse, error) {
	url := fmt.Sprintf("%s/videos/%s", c.baseURL, videoID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result VideoGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) makeRequest(ctx context.Context, endpoint string, payload interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *Client) extractVideoID(result map[string]interface{}) (string, error) {
	id, ok := result["id"].(string)
	if !ok {
		return "", ErrInvalidResponse
	}
	return id, nil
}
