package scene

import (
	"strings"
	"testing"
	"time"
)

func TestNewScene(t *testing.T) {
	chapterID := "chapter-1"
	novelID := "novel-1"
	sceneNumber := 1

	scene, err := NewScene(chapterID, novelID, sceneNumber)

	if err != nil {
		t.Errorf("NewScene() unexpected error = %v", err)
		return
	}

	if scene.ID == "" {
		t.Error("NewScene() ID should not be empty")
	}

	if scene.ChapterID != chapterID {
		t.Errorf("NewScene() ChapterID = %v, want %v", scene.ChapterID, chapterID)
	}

	if scene.NovelID != novelID {
		t.Errorf("NewScene() NovelID = %v, want %v", scene.NovelID, novelID)
	}

	if scene.SceneNumber != sceneNumber {
		t.Errorf("NewScene() SceneNumber = %v, want %v", scene.SceneNumber, sceneNumber)
	}

	if scene.Status != SceneStatusPending {
		t.Errorf("NewScene() Status = %v, want %v", scene.Status, SceneStatusPending)
	}

	if scene.Dialogues == nil {
		t.Error("NewScene() Dialogues should be initialized")
	}

	if scene.CharacterIDs == nil {
		t.Error("NewScene() CharacterIDs should be initialized")
	}

	if scene.CreatedAt.IsZero() {
		t.Error("NewScene() CreatedAt should not be zero")
	}

	if scene.UpdatedAt.IsZero() {
		t.Error("NewScene() UpdatedAt should not be zero")
	}
}

func TestScene_SetDescription(t *testing.T) {
	tests := []struct {
		name        string
		description Description
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid description with full text",
			description: Description{
				FullText: "A dark and stormy night",
			},
			wantErr: false,
		},
		{
			name: "valid description with parts",
			description: Description{
				Setting:    "forest",
				Action:     "running",
				Atmosphere: "tense",
			},
			wantErr: false,
		},
		{
			name:        "empty description",
			description: Description{},
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene := &Scene{
				UpdatedAt: time.Now().Add(-time.Hour),
			}

			oldUpdatedAt := scene.UpdatedAt
			time.Sleep(time.Millisecond)

			err := scene.SetDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetDescription() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("SetDescription() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("SetDescription() unexpected error = %v", err)
				return
			}

			if scene.Description != tt.description {
				t.Error("SetDescription() description not set correctly")
			}

			if !scene.UpdatedAt.After(oldUpdatedAt) {
				t.Error("SetDescription() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestScene_AddDialogue(t *testing.T) {
	scene := &Scene{
		Dialogues: []Dialogue{},
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	dialogue := Dialogue{
		Speaker: "Hero",
		Content: "Hello, world!",
		Emotion: "happy",
	}

	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.AddDialogue(dialogue)

	if len(scene.Dialogues) != 1 {
		t.Errorf("AddDialogue() length = %v, want 1", len(scene.Dialogues))
	}

	if scene.Dialogues[0] != dialogue {
		t.Error("AddDialogue() dialogue not added correctly")
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("AddDialogue() should update UpdatedAt timestamp")
	}
}

func TestScene_SetDialogues(t *testing.T) {
	scene := &Scene{
		Dialogues: []Dialogue{},
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	dialogues := []Dialogue{
		{Speaker: "Hero", Content: "Line 1", Emotion: "neutral"},
		{Speaker: "Villain", Content: "Line 2", Emotion: "angry"},
	}

	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.SetDialogues(dialogues)

	if len(scene.Dialogues) != len(dialogues) {
		t.Errorf("SetDialogues() length = %v, want %v", len(scene.Dialogues), len(dialogues))
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetDialogues() should update UpdatedAt timestamp")
	}
}

func TestScene_AddCharacter(t *testing.T) {
	scene := &Scene{
		CharacterIDs: []string{},
		UpdatedAt:    time.Now().Add(-time.Hour),
	}

	charID := "char-1"
	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.AddCharacter(charID)

	if len(scene.CharacterIDs) != 1 {
		t.Errorf("AddCharacter() length = %v, want 1", len(scene.CharacterIDs))
	}

	if scene.CharacterIDs[0] != charID {
		t.Errorf("AddCharacter() = %v, want %v", scene.CharacterIDs[0], charID)
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("AddCharacter() should update UpdatedAt timestamp")
	}

	oldUpdatedAt = scene.UpdatedAt
	time.Sleep(time.Millisecond)
	scene.AddCharacter(charID)

	if len(scene.CharacterIDs) != 1 {
		t.Error("AddCharacter() should not add duplicate character IDs")
	}
}

func TestScene_SetCharacters(t *testing.T) {
	scene := &Scene{
		CharacterIDs: []string{},
		UpdatedAt:    time.Now().Add(-time.Hour),
	}

	charIDs := []string{"char-1", "char-2", "char-3"}
	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.SetCharacters(charIDs)

	if len(scene.CharacterIDs) != len(charIDs) {
		t.Errorf("SetCharacters() length = %v, want %v", len(scene.CharacterIDs), len(charIDs))
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetCharacters() should update UpdatedAt timestamp")
	}
}

func TestScene_SetLocation(t *testing.T) {
	scene := &Scene{
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	location := "  Forest clearing  "
	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.SetLocation(location)

	if scene.Location != strings.TrimSpace(location) {
		t.Errorf("SetLocation() = %v, want %v", scene.Location, strings.TrimSpace(location))
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetLocation() should update UpdatedAt timestamp")
	}
}

func TestScene_SetTimeOfDay(t *testing.T) {
	scene := &Scene{
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	timeOfDay := "  Dawn  "
	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.SetTimeOfDay(timeOfDay)

	if scene.TimeOfDay != strings.TrimSpace(timeOfDay) {
		t.Errorf("SetTimeOfDay() = %v, want %v", scene.TimeOfDay, strings.TrimSpace(timeOfDay))
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetTimeOfDay() should update UpdatedAt timestamp")
	}
}

func TestScene_SetImagePrompt(t *testing.T) {
	scene := &Scene{
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	prompt := "a beautiful landscape"
	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.SetImagePrompt(prompt)

	if scene.ImagePrompt != prompt {
		t.Errorf("SetImagePrompt() = %v, want %v", scene.ImagePrompt, prompt)
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetImagePrompt() should update UpdatedAt timestamp")
	}
}

func TestScene_SetVideoPrompt(t *testing.T) {
	scene := &Scene{
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	prompt := "a moving scene"
	oldUpdatedAt := scene.UpdatedAt
	time.Sleep(time.Millisecond)

	scene.SetVideoPrompt(prompt)

	if scene.VideoPrompt != prompt {
		t.Errorf("SetVideoPrompt() = %v, want %v", scene.VideoPrompt, prompt)
	}

	if !scene.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetVideoPrompt() should update UpdatedAt timestamp")
	}
}

func TestScene_UpdateStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      SceneStatus
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "valid status - pending",
			status:  SceneStatusPending,
			wantErr: false,
		},
		{
			name:    "valid status - ready",
			status:  SceneStatusReady,
			wantErr: false,
		},
		{
			name:    "valid status - generating",
			status:  SceneStatusGenerating,
			wantErr: false,
		},
		{
			name:    "valid status - completed",
			status:  SceneStatusCompleted,
			wantErr: false,
		},
		{
			name:    "valid status - failed",
			status:  SceneStatusFailed,
			wantErr: false,
		},
		{
			name:        "invalid status",
			status:      SceneStatus("invalid"),
			wantErr:     true,
			expectedErr: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene := &Scene{
				UpdatedAt: time.Now().Add(-time.Hour),
			}

			oldUpdatedAt := scene.UpdatedAt
			time.Sleep(time.Millisecond)

			err := scene.UpdateStatus(tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateStatus() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("UpdateStatus() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStatus() unexpected error = %v", err)
				return
			}

			if scene.Status != tt.status {
				t.Errorf("UpdateStatus() Status = %v, want %v", scene.Status, tt.status)
			}

			if !scene.UpdatedAt.After(oldUpdatedAt) {
				t.Error("UpdateStatus() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestScene_Validate(t *testing.T) {
	tests := []struct {
		name        string
		scene       *Scene
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid scene",
			scene: &Scene{
				Description: Description{
					FullText: "A scene description",
				},
			},
			wantErr: false,
		},
		{
			name: "valid scene with description parts",
			scene: &Scene{
				Description: Description{
					Setting: "forest",
					Action:  "running",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid scene - empty description",
			scene: &Scene{
				Description: Description{},
			},
			wantErr:     true,
			expectedErr: ErrEmptyDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.scene.Validate()

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

func TestScene_HasCharacters(t *testing.T) {
	tests := []struct {
		name  string
		scene *Scene
		want  bool
	}{
		{
			name: "has characters",
			scene: &Scene{
				CharacterIDs: []string{"char-1", "char-2"},
			},
			want: true,
		},
		{
			name: "no characters",
			scene: &Scene{
				CharacterIDs: []string{},
			},
			want: false,
		},
		{
			name: "nil characters",
			scene: &Scene{
				CharacterIDs: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scene.HasCharacters(); got != tt.want {
				t.Errorf("HasCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_HasDialogues(t *testing.T) {
	tests := []struct {
		name  string
		scene *Scene
		want  bool
	}{
		{
			name: "has dialogues",
			scene: &Scene{
				Dialogues: []Dialogue{{Speaker: "Hero", Content: "Hi"}},
			},
			want: true,
		},
		{
			name: "no dialogues",
			scene: &Scene{
				Dialogues: []Dialogue{},
			},
			want: false,
		},
		{
			name: "nil dialogues",
			scene: &Scene{
				Dialogues: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scene.HasDialogues(); got != tt.want {
				t.Errorf("HasDialogues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescription_IsEmpty(t *testing.T) {
	tests := []struct {
		name        string
		description Description
		want        bool
	}{
		{
			name:        "empty description",
			description: Description{},
			want:        true,
		},
		{
			name: "has full text",
			description: Description{
				FullText: "description",
			},
			want: false,
		},
		{
			name: "has setting",
			description: Description{
				Setting: "forest",
			},
			want: false,
		},
		{
			name: "has action",
			description: Description{
				Action: "running",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.description.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescription_ToPrompt(t *testing.T) {
	tests := []struct {
		name        string
		description Description
		wantEmpty   bool
		contains    string
	}{
		{
			name: "full text only",
			description: Description{
				FullText: "A complete description",
			},
			wantEmpty: false,
			contains:  "A complete description",
		},
		{
			name: "parts without full text",
			description: Description{
				Setting:    "forest",
				Action:     "running",
				Atmosphere: "tense",
			},
			wantEmpty: false,
			contains:  "forest",
		},
		{
			name:        "empty description",
			description: Description{},
			wantEmpty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.description.ToPrompt()
			if tt.wantEmpty && got != "" {
				t.Errorf("ToPrompt() = %v, want empty string", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Error("ToPrompt() should not return empty string")
			}
			if tt.contains != "" && !strings.Contains(got, tt.contains) {
				t.Errorf("ToPrompt() = %v, should contain %v", got, tt.contains)
			}
		})
	}
}
