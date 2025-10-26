package character

import (
	"strings"
	"testing"
	"time"
)

func TestNewCharacter(t *testing.T) {
	tests := []struct {
		name        string
		novelID     string
		charName    string
		role        CharacterRole
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "valid main character",
			novelID:  "novel-1",
			charName: "Test Character",
			role:     CharacterRoleMain,
			wantErr:  false,
		},
		{
			name:     "valid supporting character",
			novelID:  "novel-1",
			charName: "Supporting Char",
			role:     CharacterRoleSupporting,
			wantErr:  false,
		},
		{
			name:     "valid minor character",
			novelID:  "novel-1",
			charName: "Minor Char",
			role:     CharacterRoleMinor,
			wantErr:  false,
		},
		{
			name:        "empty name",
			novelID:     "novel-1",
			charName:    "",
			role:        CharacterRoleMain,
			wantErr:     true,
			expectedErr: ErrEmptyCharacterName,
		},
		{
			name:        "whitespace name",
			novelID:     "novel-1",
			charName:    "   ",
			role:        CharacterRoleMain,
			wantErr:     true,
			expectedErr: ErrEmptyCharacterName,
		},
		{
			name:        "invalid role",
			novelID:     "novel-1",
			charName:    "Test Character",
			role:        CharacterRole("invalid"),
			wantErr:     true,
			expectedErr: ErrInvalidRole,
		},
		{
			name:     "name with whitespace trimmed",
			novelID:  "novel-1",
			charName: "  Test Character  ",
			role:     CharacterRoleMain,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, err := NewCharacter(tt.novelID, tt.charName, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewCharacter() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("NewCharacter() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewCharacter() unexpected error = %v", err)
				return
			}

			if char.ID == "" {
				t.Error("NewCharacter() ID should not be empty")
			}

			if char.NovelID != tt.novelID {
				t.Errorf("NewCharacter() NovelID = %v, want %v", char.NovelID, tt.novelID)
			}

			if char.Name != strings.TrimSpace(tt.charName) {
				t.Errorf("NewCharacter() Name = %v, want %v", char.Name, strings.TrimSpace(tt.charName))
			}

			if char.Role != tt.role {
				t.Errorf("NewCharacter() Role = %v, want %v", char.Role, tt.role)
			}

			if char.CreatedAt.IsZero() {
				t.Error("NewCharacter() CreatedAt should not be zero")
			}

			if char.UpdatedAt.IsZero() {
				t.Error("NewCharacter() UpdatedAt should not be zero")
			}
		})
	}
}

func TestCharacter_SetAppearance(t *testing.T) {
	tests := []struct {
		name        string
		appearance  Appearance
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid appearance",
			appearance: Appearance{
				PhysicalTraits:   "tall, blonde",
				ClothingStyle:    "casual",
				DistinctFeatures: "scar on left cheek",
			},
			wantErr: false,
		},
		{
			name: "appearance with partial fields",
			appearance: Appearance{
				PhysicalTraits: "short hair",
			},
			wantErr: false,
		},
		{
			name:        "empty appearance",
			appearance:  Appearance{},
			wantErr:     true,
			expectedErr: ErrEmptyAppearance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &Character{
				Name:      "Test",
				UpdatedAt: time.Now().Add(-time.Hour),
			}

			oldUpdatedAt := char.UpdatedAt
			time.Sleep(time.Millisecond)

			err := char.SetAppearance(tt.appearance)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetAppearance() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("SetAppearance() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("SetAppearance() unexpected error = %v", err)
				return
			}

			if char.Appearance != tt.appearance {
				t.Error("SetAppearance() appearance not set correctly")
			}

			if !char.UpdatedAt.After(oldUpdatedAt) {
				t.Error("SetAppearance() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestCharacter_SetPersonality(t *testing.T) {
	char := &Character{
		Name:      "Test",
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	personality := Personality{
		Traits:     "brave, loyal",
		Motivation: "save the world",
		Background: "orphan raised by monks",
	}

	oldUpdatedAt := char.UpdatedAt
	time.Sleep(time.Millisecond)

	char.SetPersonality(personality)

	if char.Personality != personality {
		t.Error("SetPersonality() personality not set correctly")
	}

	if !char.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetPersonality() should update UpdatedAt timestamp")
	}
}

func TestCharacter_SetDescription(t *testing.T) {
	char := &Character{
		Name:      "Test",
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	description := "  A brave warrior  "
	oldUpdatedAt := char.UpdatedAt
	time.Sleep(time.Millisecond)

	char.SetDescription(description)

	if char.Description != strings.TrimSpace(description) {
		t.Errorf("SetDescription() = %v, want %v", char.Description, strings.TrimSpace(description))
	}

	if !char.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetDescription() should update UpdatedAt timestamp")
	}
}

func TestCharacter_SetReferenceImage(t *testing.T) {
	char := &Character{
		Name:      "Test",
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	imageURL := "https://example.com/image.jpg"
	oldUpdatedAt := char.UpdatedAt
	time.Sleep(time.Millisecond)

	char.SetReferenceImage(imageURL)

	if char.ReferenceImageURL != imageURL {
		t.Errorf("SetReferenceImage() = %v, want %v", char.ReferenceImageURL, imageURL)
	}

	if !char.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetReferenceImage() should update UpdatedAt timestamp")
	}

	if !char.HasReferenceImage() {
		t.Error("HasReferenceImage() should return true after setting image")
	}
}

func TestCharacter_HasReferenceImage(t *testing.T) {
	tests := []struct {
		name  string
		char  *Character
		want  bool
	}{
		{
			name: "has reference image",
			char: &Character{
				ReferenceImageURL: "https://example.com/image.jpg",
			},
			want: true,
		},
		{
			name: "no reference image",
			char: &Character{
				ReferenceImageURL: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.char.HasReferenceImage(); got != tt.want {
				t.Errorf("HasReferenceImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_Validate(t *testing.T) {
	tests := []struct {
		name        string
		char        *Character
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid character",
			char: &Character{
				Name: "Test Character",
				Role: CharacterRoleMain,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			char: &Character{
				Name: "",
				Role: CharacterRoleMain,
			},
			wantErr:     true,
			expectedErr: ErrEmptyCharacterName,
		},
		{
			name: "invalid role",
			char: &Character{
				Name: "Test Character",
				Role: CharacterRole("invalid"),
			},
			wantErr:     true,
			expectedErr: ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.char.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("Validate() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}

func TestAppearance_IsEmpty(t *testing.T) {
	tests := []struct {
		name       string
		appearance Appearance
		want       bool
	}{
		{
			name:       "empty appearance",
			appearance: Appearance{},
			want:       true,
		},
		{
			name: "has physical traits",
			appearance: Appearance{
				PhysicalTraits: "tall",
			},
			want: false,
		},
		{
			name: "has clothing style",
			appearance: Appearance{
				ClothingStyle: "casual",
			},
			want: false,
		},
		{
			name: "has distinct features",
			appearance: Appearance{
				DistinctFeatures: "scar",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appearance.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppearance_ToPrompt(t *testing.T) {
	tests := []struct {
		name       string
		appearance Appearance
		wantEmpty  bool
	}{
		{
			name: "full appearance",
			appearance: Appearance{
				Age:              "25 years old",
				Height:           "tall",
				PhysicalTraits:   "blonde hair",
				ClothingStyle:    "casual",
				DistinctFeatures: "scar",
			},
			wantEmpty: false,
		},
		{
			name:       "empty appearance",
			appearance: Appearance{},
			wantEmpty:  true,
		},
		{
			name: "partial appearance",
			appearance: Appearance{
				PhysicalTraits: "red hair",
				Age:            "young",
			},
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.appearance.ToPrompt()
			if tt.wantEmpty && got != "" {
				t.Errorf("ToPrompt() = %v, want empty string", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Error("ToPrompt() should not return empty string")
			}
		})
	}
}

func TestPersonality_IsEmpty(t *testing.T) {
	tests := []struct {
		name        string
		personality Personality
		want        bool
	}{
		{
			name:        "empty personality",
			personality: Personality{},
			want:        true,
		},
		{
			name: "has traits",
			personality: Personality{
				Traits: "brave",
			},
			want: false,
		},
		{
			name: "has motivation",
			personality: Personality{
				Motivation: "save the world",
			},
			want: false,
		},
		{
			name: "has background",
			personality: Personality{
				Background: "orphan",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.personality.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_GeneratePrompt(t *testing.T) {
	tests := []struct {
		name      string
		char      *Character
		wantEmpty bool
	}{
		{
			name: "character with full details",
			char: &Character{
				Name: "Hero",
				Appearance: Appearance{
					PhysicalTraits: "tall, blonde",
				},
				Description: "brave warrior",
			},
			wantEmpty: false,
		},
		{
			name: "character with name only",
			char: &Character{
				Name: "Hero",
			},
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.char.GeneratePrompt()
			if tt.wantEmpty && got != "" {
				t.Errorf("GeneratePrompt() = %v, want empty", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Error("GeneratePrompt() should not return empty string")
			}
			if !strings.Contains(got, tt.char.Name) {
				t.Errorf("GeneratePrompt() should contain character name: %v", tt.char.Name)
			}
		})
	}
}
