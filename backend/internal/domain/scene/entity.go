package scene

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SceneID string

type SceneStatus string

const (
	SceneStatusPending    SceneStatus = "pending"
	SceneStatusReady      SceneStatus = "ready"
	SceneStatusGenerating SceneStatus = "generating"
	SceneStatusCompleted  SceneStatus = "completed"
	SceneStatusFailed     SceneStatus = "failed"
)

var (
	ErrSceneNotFound    = errors.New("scene not found")
	ErrInvalidScene     = errors.New("invalid scene")
	ErrEmptyDescription = errors.New("scene description cannot be empty")
	ErrInvalidStatus    = errors.New("invalid scene status")
)

type Scene struct {
	ID           SceneID
	ChapterID    string
	NovelID      string
	SceneNumber  int
	Location     string
	TimeOfDay    string
	Description  Description
	Dialogues    []Dialogue
	CharacterIDs []string
	ImagePrompt  string
	VideoPrompt  string
	Status       SceneStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Description struct {
	Setting    string
	Action     string
	Atmosphere string
	FullText   string
}

func (d Description) IsEmpty() bool {
	return d.FullText == "" && d.Setting == "" && d.Action == ""
}

func (d Description) ToPrompt() string {
	if d.FullText != "" {
		return d.FullText
	}

	var parts []string
	if d.Setting != "" {
		parts = append(parts, d.Setting)
	}
	if d.Action != "" {
		parts = append(parts, d.Action)
	}
	if d.Atmosphere != "" {
		parts = append(parts, d.Atmosphere)
	}

	return strings.Join(parts, ". ")
}

type Dialogue struct {
	Speaker string
	Content string
	Emotion string
}

func NewScene(chapterID, novelID string, sceneNumber int) (*Scene, error) {
	now := time.Now()
	return &Scene{
		ID:           SceneID(uuid.New().String()),
		ChapterID:    chapterID,
		NovelID:      novelID,
		SceneNumber:  sceneNumber,
		Status:       SceneStatusPending,
		Dialogues:    []Dialogue{},
		CharacterIDs: []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (s *Scene) SetDescription(description Description) error {
	if description.IsEmpty() {
		return ErrEmptyDescription
	}
	s.Description = description
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Scene) AddDialogue(dialogue Dialogue) {
	s.Dialogues = append(s.Dialogues, dialogue)
	s.UpdatedAt = time.Now()
}

func (s *Scene) SetDialogues(dialogues []Dialogue) {
	s.Dialogues = dialogues
	s.UpdatedAt = time.Now()
}

func (s *Scene) AddCharacter(characterID string) {
	for _, id := range s.CharacterIDs {
		if id == characterID {
			return
		}
	}
	s.CharacterIDs = append(s.CharacterIDs, characterID)
	s.UpdatedAt = time.Now()
}

func (s *Scene) SetCharacters(characterIDs []string) {
	s.CharacterIDs = characterIDs
	s.UpdatedAt = time.Now()
}

func (s *Scene) SetLocation(location string) {
	s.Location = strings.TrimSpace(location)
	s.UpdatedAt = time.Now()
}

func (s *Scene) SetTimeOfDay(timeOfDay string) {
	s.TimeOfDay = strings.TrimSpace(timeOfDay)
	s.UpdatedAt = time.Now()
}

func (s *Scene) SetImagePrompt(prompt string) {
	s.ImagePrompt = prompt
	s.UpdatedAt = time.Now()
}

func (s *Scene) SetVideoPrompt(prompt string) {
	s.VideoPrompt = prompt
	s.UpdatedAt = time.Now()
}

func (s *Scene) UpdateStatus(status SceneStatus) error {
	validStatuses := map[SceneStatus]bool{
		SceneStatusPending:    true,
		SceneStatusReady:      true,
		SceneStatusGenerating: true,
		SceneStatusCompleted:  true,
		SceneStatusFailed:     true,
	}

	if !validStatuses[status] {
		return ErrInvalidStatus
	}

	s.Status = status
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Scene) Validate() error {
	if s.Description.IsEmpty() {
		return ErrEmptyDescription
	}
	return nil
}

func (s *Scene) HasCharacters() bool {
	return len(s.CharacterIDs) > 0
}

func (s *Scene) HasDialogues() bool {
	return len(s.Dialogues) > 0
}
