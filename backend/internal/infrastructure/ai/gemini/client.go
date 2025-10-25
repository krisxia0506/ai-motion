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
	ErrInvalidResponse = errors.New("invalid response from gemini api")
	ErrRateLimited     = errors.New("rate limited by gemini api")
	ErrContentFiltered = errors.New("content filtered by gemini api")
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, ErrAPIKeyRequired
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
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

type GenerateResponse struct {
	Images []struct {
		URL           string `json:"url"`
		B64JSON       string `json:"b64_json"`
		RevisedPrompt string `json:"revised_prompt"`
	} `json:"images"`
	Created int64 `json:"created"`
}

func (c *Client) TextToImage(ctx context.Context, req TextToImageRequest) (string, error) {
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": req.Prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.7,
		},
	}

	if req.NegativePrompt != "" {
		payload["safetySettings"] = []map[string]interface{}{
			{
				"category":  "HARM_CATEGORY_HARASSMENT",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
		}
	}

	result, err := c.makeRequest(ctx, "models/gemini-2.0-flash-exp:generateContent", payload)
	if err != nil {
		return "", err
	}

	return c.extractImageURL(result)
}

func (c *Client) ImageToImage(ctx context.Context, req ImageToImageRequest) (string, error) {
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": req.Prompt},
					{
						"inline_data": map[string]string{
							"mime_type": "image/jpeg",
							"data":      req.ReferenceImage,
						},
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.7,
		},
	}

	result, err := c.makeRequest(ctx, "models/gemini-2.0-flash-exp:generateContent", payload)
	if err != nil {
		return "", err
	}

	return c.extractImageURL(result)
}

func (c *Client) makeRequest(ctx context.Context, endpoint string, payload interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s?key=%s", c.baseURL, endpoint, c.apiKey)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", ErrInvalidResponse
	}

	candidate := candidates[0].(map[string]interface{})
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return "", ErrInvalidResponse
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", ErrInvalidResponse
	}

	part := parts[0].(map[string]interface{})
	if inlineData, ok := part["inline_data"].(map[string]interface{}); ok {
		if data, ok := inlineData["data"].(string); ok {
			return data, nil
		}
	}

	return "", ErrInvalidResponse
}
