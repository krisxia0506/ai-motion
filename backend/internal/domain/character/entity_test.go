package character_test

import (
	"testing"

	"github.com/xiajiayi/ai-motion/internal/domain/character"
)

func TestNewCharacter(t *testing.T) {
	tests := []struct {
		name     string
		novelID  string
		charName string
		role     character.CharacterRole
		wantErr  bool
	}{
		{
			name:     "valid main character",
			novelID:  "novel-1",
			charName: "张三",
			role:     character.CharacterRoleMain,
			wantErr:  false,
		},
		{
			name:     "valid supporting character",
			novelID:  "novel-1",
			charName: "李四",
			role:     character.CharacterRoleSupporting,
			wantErr:  false,
		},
		{
			name:     "empty name",
			novelID:  "novel-1",
			charName: "",
			role:     character.CharacterRoleMain,
			wantErr:  true,
		},
		{
			name:     "invalid role",
			novelID:  "novel-1",
			charName: "王五",
			role:     character.CharacterRole("invalid"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, err := character.NewCharacter(tt.novelID, tt.charName, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewCharacter() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewCharacter() unexpected error: %v", err)
			}

			if char.Name != tt.charName {
				t.Errorf("Name = %v, want %v", char.Name, tt.charName)
			}

			if char.Role != tt.role {
				t.Errorf("Role = %v, want %v", char.Role, tt.role)
			}

			if char.NovelID != tt.novelID {
				t.Errorf("NovelID = %v, want %v", char.NovelID, tt.novelID)
			}
		})
	}
}

func TestCharacter_SetAppearance(t *testing.T) {
	char := &character.Character{
		Name: "Test Character",
		Role: character.CharacterRoleMain,
	}

	tests := []struct {
		name       string
		appearance character.Appearance
		wantErr    bool
	}{
		{
			name: "valid appearance",
			appearance: character.Appearance{
				PhysicalTraits: "高大威猛",
				ClothingStyle:  "长袍",
			},
			wantErr: false,
		},
		{
			name:       "empty appearance",
			appearance: character.Appearance{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := char.SetAppearance(tt.appearance)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetAppearance() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("SetAppearance() unexpected error: %v", err)
			}

			if char.Appearance.IsEmpty() {
				t.Errorf("Appearance should not be empty")
			}
		})
	}
}

func TestAppearance_ToPrompt(t *testing.T) {
	tests := []struct {
		name       string
		appearance character.Appearance
		want       string
	}{
		{
			name: "full appearance",
			appearance: character.Appearance{
				Age:              "25岁",
				Height:           "180cm",
				PhysicalTraits:   "黑发黑眼",
				ClothingStyle:    "白色长袍",
				DistinctFeatures: "眉心有痣",
			},
			want: "25岁, 180cm, 黑发黑眼, 白色长袍, 眉心有痣",
		},
		{
			name: "partial appearance",
			appearance: character.Appearance{
				PhysicalTraits: "黑发",
			},
			want: "黑发",
		},
		{
			name:       "empty appearance",
			appearance: character.Appearance{},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.appearance.ToPrompt()
			if got != tt.want {
				t.Errorf("ToPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_GeneratePrompt(t *testing.T) {
	char := &character.Character{
		Name: "张三",
		Appearance: character.Appearance{
			PhysicalTraits: "黑发黑眼",
		},
		Description: "主角",
	}

	prompt := char.GeneratePrompt()

	if prompt == "" {
		t.Errorf("GeneratePrompt() should not be empty")
	}

	if len(prompt) < 3 {
		t.Errorf("GeneratePrompt() should contain character information")
	}
}

func TestCharacter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		char    *character.Character
		wantErr bool
	}{
		{
			name: "valid character",
			char: &character.Character{
				Name: "张三",
				Role: character.CharacterRoleMain,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			char: &character.Character{
				Name: "",
				Role: character.CharacterRoleMain,
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			char: &character.Character{
				Name: "张三",
				Role: character.CharacterRole("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.char.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("Validate() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}
