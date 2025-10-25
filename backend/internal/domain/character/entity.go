package character

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CharacterID string

type CharacterRole string

const (
	CharacterRoleMain       CharacterRole = "main"
	CharacterRoleSupporting CharacterRole = "supporting"
	CharacterRoleMinor      CharacterRole = "minor"
)

var (
	ErrCharacterNotFound  = errors.New("character not found")
	ErrInvalidCharacter   = errors.New("invalid character")
	ErrEmptyCharacterName = errors.New("character name cannot be empty")
	ErrInvalidRole        = errors.New("invalid character role")
	ErrEmptyAppearance    = errors.New("character appearance cannot be empty")
)

type Character struct {
	ID                CharacterID
	NovelID           string
	Name              string
	Role              CharacterRole
	Appearance        Appearance
	Personality       Personality
	Description       string
	ReferenceImageURL string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Appearance struct {
	PhysicalTraits   string
	ClothingStyle    string
	DistinctFeatures string
	Age              string
	Height           string
}

func (a Appearance) IsEmpty() bool {
	return a.PhysicalTraits == "" && a.ClothingStyle == "" && a.DistinctFeatures == ""
}

func (a Appearance) ToPrompt() string {
	var parts []string

	if a.Age != "" {
		parts = append(parts, a.Age)
	}
	if a.Height != "" {
		parts = append(parts, a.Height)
	}
	if a.PhysicalTraits != "" {
		parts = append(parts, a.PhysicalTraits)
	}
	if a.ClothingStyle != "" {
		parts = append(parts, a.ClothingStyle)
	}
	if a.DistinctFeatures != "" {
		parts = append(parts, a.DistinctFeatures)
	}

	return strings.Join(parts, ", ")
}

type Personality struct {
	Traits     string
	Motivation string
	Background string
}

func (p Personality) IsEmpty() bool {
	return p.Traits == "" && p.Motivation == "" && p.Background == ""
}

func NewCharacter(novelID, name string, role CharacterRole) (*Character, error) {
	if err := validateCharacterInput(name, role); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Character{
		ID:        CharacterID(uuid.New().String()),
		NovelID:   novelID,
		Name:      strings.TrimSpace(name),
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Character) SetAppearance(appearance Appearance) error {
	if appearance.IsEmpty() {
		return ErrEmptyAppearance
	}
	c.Appearance = appearance
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Character) SetPersonality(personality Personality) {
	c.Personality = personality
	c.UpdatedAt = time.Now()
}

func (c *Character) SetDescription(description string) {
	c.Description = strings.TrimSpace(description)
	c.UpdatedAt = time.Now()
}

func (c *Character) SetReferenceImage(imageURL string) {
	c.ReferenceImageURL = imageURL
	c.UpdatedAt = time.Now()
}

func (c *Character) HasReferenceImage() bool {
	return c.ReferenceImageURL != ""
}

func (c *Character) Validate() error {
	if c.Name == "" {
		return ErrEmptyCharacterName
	}
	if !isValidRole(c.Role) {
		return ErrInvalidRole
	}
	return nil
}

func validateCharacterInput(name string, role CharacterRole) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyCharacterName
	}
	if !isValidRole(role) {
		return ErrInvalidRole
	}
	return nil
}

func isValidRole(role CharacterRole) bool {
	return role == CharacterRoleMain ||
		role == CharacterRoleSupporting ||
		role == CharacterRoleMinor
}

func (c *Character) GeneratePrompt() string {
	var parts []string

	parts = append(parts, c.Name)

	if !c.Appearance.IsEmpty() {
		parts = append(parts, c.Appearance.ToPrompt())
	}

	if c.Description != "" {
		parts = append(parts, c.Description)
	}

	return strings.Join(parts, ", ")
}
