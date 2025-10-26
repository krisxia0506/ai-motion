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

	client := postgrest.NewClient(cfg.URL, cfg.APIKey, nil)
	if client.ClientError != nil {
		return nil, fmt.Errorf("failed to create supabase client: %w", client.ClientError)
	}

	return client, nil
}
