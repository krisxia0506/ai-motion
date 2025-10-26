package database

import (
	"fmt"

	postgrest "github.com/supabase-community/postgrest-go"
)

type SupabaseConfig struct {
	URL    string
	APIKey string
}

func NewSupabaseClient(cfg *SupabaseConfig) (*postgrest.Client, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("supabase URL is required")
	}

	// NewClient takes (url, schema, headers)
	// We need to set the API key via SetApiKey method
	client := postgrest.NewClient(cfg.URL, "", nil)
	if client.ClientError != nil {
		return nil, fmt.Errorf("failed to create supabase client: %w", client.ClientError)
	}

	// Set the service role API key
	client = client.SetApiKey(cfg.APIKey)

	return client, nil
}
