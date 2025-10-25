package scene_test

import (
	"testing"

	"github.com/xiajiayi/ai-motion/internal/domain/scene"
)

func TestNewScene(t *testing.T) {
	chapterID := "chapter-1"
	novelID := "novel-1"
	sceneNumber := 1

	s, err := scene.NewScene(chapterID, novelID, sceneNumber)

	if err != nil {
		t.Errorf("NewScene() unexpected error: %v", err)
	}

	if s.ChapterID != chapterID {
		t.Errorf("ChapterID = %v, want %v", s.ChapterID, chapterID)
	}

	if s.NovelID != novelID {
		t.Errorf("NovelID = %v, want %v", s.NovelID, novelID)
	}

	if s.SceneNumber != sceneNumber {
		t.Errorf("SceneNumber = %v, want %v", s.SceneNumber, sceneNumber)
	}

	if s.Status != scene.SceneStatusPending {
		t.Errorf("Status = %v, want %v", s.Status, scene.SceneStatusPending)
	}
}

func TestScene_SetDescription(t *testing.T) {
	s := &scene.Scene{
		ChapterID: "chapter-1",
	}

	tests := []struct {
		name        string
		description scene.Description
		wantErr     bool
	}{
		{
			name: "valid description",
			description: scene.Description{
				FullText: "这是一个测试场景描述",
			},
			wantErr: false,
		},
		{
			name: "description with components",
			description: scene.Description{
				Setting:    "室内",
				Action:     "说话",
				Atmosphere: "紧张",
			},
			wantErr: false,
		},
		{
			name:        "empty description",
			description: scene.Description{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.SetDescription(tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetDescription() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("SetDescription() unexpected error: %v", err)
			}
		})
	}
}

func TestDescription_ToPrompt(t *testing.T) {
	tests := []struct {
		name string
		desc scene.Description
		want string
	}{
		{
			name: "full text only",
			desc: scene.Description{
				FullText: "完整描述",
			},
			want: "完整描述",
		},
		{
			name: "components only",
			desc: scene.Description{
				Setting:    "室内",
				Action:     "说话",
				Atmosphere: "紧张",
			},
			want: "室内. 说话. 紧张",
		},
		{
			name: "empty description",
			desc: scene.Description{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.desc.ToPrompt()
			if got != tt.want {
				t.Errorf("ToPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_AddDialogue(t *testing.T) {
	s := &scene.Scene{
		ChapterID: "chapter-1",
	}

	dialogue := scene.Dialogue{
		Speaker: "张三",
		Content: "你好",
		Emotion: "happy",
	}

	s.AddDialogue(dialogue)

	if len(s.Dialogues) != 1 {
		t.Errorf("len(Dialogues) = %v, want 1", len(s.Dialogues))
	}

	if s.Dialogues[0].Speaker != "张三" {
		t.Errorf("Dialogue speaker = %v, want 张三", s.Dialogues[0].Speaker)
	}
}

func TestScene_AddCharacter(t *testing.T) {
	s := &scene.Scene{
		ChapterID: "chapter-1",
	}

	s.AddCharacter("char-1")
	s.AddCharacter("char-2")
	s.AddCharacter("char-1")

	if len(s.CharacterIDs) != 2 {
		t.Errorf("len(CharacterIDs) = %v, want 2", len(s.CharacterIDs))
	}
}

func TestScene_UpdateStatus(t *testing.T) {
	s := &scene.Scene{
		ChapterID: "chapter-1",
		Status:    scene.SceneStatusPending,
	}

	tests := []struct {
		name    string
		status  scene.SceneStatus
		wantErr bool
	}{
		{
			name:    "valid status - ready",
			status:  scene.SceneStatusReady,
			wantErr: false,
		},
		{
			name:    "valid status - generating",
			status:  scene.SceneStatusGenerating,
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  scene.SceneStatus("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdateStatus(tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateStatus() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStatus() unexpected error: %v", err)
			}

			if s.Status != tt.status {
				t.Errorf("Status = %v, want %v", s.Status, tt.status)
			}
		})
	}
}

func TestScene_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scene   *scene.Scene
		wantErr bool
	}{
		{
			name: "valid scene",
			scene: &scene.Scene{
				Description: scene.Description{
					FullText: "测试描述",
				},
			},
			wantErr: false,
		},
		{
			name: "empty description",
			scene: &scene.Scene{
				Description: scene.Description{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.scene.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("Validate() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}
