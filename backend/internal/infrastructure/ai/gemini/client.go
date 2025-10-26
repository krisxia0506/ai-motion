package gemini

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
	ErrAPIKeyRequired  = errors.New("gemini api key is required")
	ErrBaseURLRequired = errors.New("gemini base url is required")
	ErrInvalidResponse = errors.New("invalid response from gemini api")
	ErrRateLimited     = errors.New("rate limited by gemini api")
	ErrContentFiltered = errors.New("content filtered by gemini api")
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
			Timeout: 60 * time.Second,
		},
	}, nil
}

type TextToImageRequest struct {
	Prompt         string
	NegativePrompt string
	Width          int
	Height         int
	Quality        string
	Style          string
}

type ImageToImageRequest struct {
	ReferenceImage string
	Prompt         string
	NegativePrompt string
	Strength       float64
	Width          int
	Height         int
}

type OpenAIImageResponse struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL           string `json:"url"`
		B64JSON       string `json:"b64_json"`
		RevisedPrompt string `json:"revised_prompt"`
	} `json:"data"`
}

func (c *Client) TextToImage(ctx context.Context, req TextToImageRequest) (string, error) {
	// 使用固定尺寸 1344x768
	size := "1344x768"

	payload := map[string]interface{}{
		"model":  "gemini-2.5-flash-image",
		"prompt": req.Prompt,
		"n":      1,
		"size":   size,
	}

	result, err := c.makeRequest(ctx, "images/generations", payload)
	if err != nil {
		return "", err
	}

	return c.extractImageURL(result)
}

func (c *Client) ImageToImage(ctx context.Context, req ImageToImageRequest) (string, error) {
	size := "1344x768"

	payload := map[string]interface{}{
		"model":  "gemini-2.5-flash-image",
		"prompt": req.Prompt,
		"image":  req.ReferenceImage,
		"n":      1,
		"size":   size,
	}

	if req.Strength > 0 {
		payload["strength"] = req.Strength
	}

	result, err := c.makeRequest(ctx, "images/generations", payload)
	if err != nil {
		return "", err
	}

	return c.extractImageURL(result)
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

func (c *Client) extractImageURL(result map[string]interface{}) (string, error) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return "", ErrInvalidResponse
	}

	firstImage := data[0].(map[string]interface{})

	if url, ok := firstImage["url"].(string); ok && url != "" {
		return url, nil
	}

	if b64JSON, ok := firstImage["b64_json"].(string); ok && b64JSON != "" {
		return b64JSON, nil
	}

	return "", ErrInvalidResponse
}
