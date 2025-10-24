---
name: implement-ddd-entity
description: Complete DDD entity implementation guide following clean architecture
version: 1.0.0
author: AI-Motion Team
---

# Implement DDD Entity Skill

You are an expert in implementing Domain-Driven Design entities for the AI-Motion project.

## DDD Entity Implementation Guide

When implementing a new DDD entity, follow these steps strictly:

## Step 1: Define Domain Entity

**Location:** `backend/internal/domain/{entity}/`

Create the following files:

### 1.1 Entity Definition (`entity.go`)

\`\`\`go
// Package {entity} implements the {Entity} aggregate root and related domain logic.
package {entity}

import (
    "time"
    "errors"
)

// {Entity}ID represents a unique identifier for {Entity}.
type {Entity}ID string

// {Entity}Status represents the lifecycle status of a {Entity}.
type {Entity}Status string

const (
    {Entity}StatusCreated   {Entity}Status = "created"
    {Entity}StatusProcessing {Entity}Status = "processing"
    {Entity}StatusCompleted  {Entity}Status = "completed"
    {Entity}StatusFailed     {Entity}Status = "failed"
)

// {Entity} represents the {entity} aggregate root.
// It encapsulates all business rules and invariants for {entity} management.
type {Entity} struct {
    ID        {Entity}ID
    // Add fields
    Status    {Entity}Status
    CreatedAt time.Time
    UpdatedAt time.Time
}

// New{Entity} creates a new {Entity} instance with required fields.
// It enforces domain invariants at creation time.
func New{Entity}(/* params */) *{Entity} {
    return &{Entity}{
        ID:        {Entity}ID(generateID()),
        // Initialize fields
        Status:    {Entity}StatusCreated,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// Business logic methods

// Validate checks if the {Entity} satisfies all domain invariants.
func (e *{Entity}) Validate() error {
    if e.ID == "" {
        return errors.New("id cannot be empty")
    }
    // Add validation logic
    return nil
}

// Add domain behavior methods here
// Example: Process, Complete, Fail, etc.
\`\`\`

### 1.2 Value Objects (`value_object.go`)

\`\`\`go
package {entity}

// Define value objects related to this entity
// Value objects are immutable and compared by value

// Example{ValueObject} represents a value object in the {entity} domain.
type Example{ValueObject} struct {
    Field1 string
    Field2 int
}

// NewExample{ValueObject} creates a new value object with validation.
func NewExample{ValueObject}(field1 string, field2 int) (Example{ValueObject}, error) {
    vo := Example{ValueObject}{
        Field1: field1,
        Field2: field2,
    }

    if err := vo.validate(); err != nil {
        return Example{ValueObject}{}, err
    }

    return vo, nil
}

func (v Example{ValueObject}) validate() error {
    // Validation logic
    return nil
}

// Value objects should have equality methods
func (v Example{ValueObject}) Equals(other Example{ValueObject}) bool {
    return v.Field1 == other.Field1 && v.Field2 == other.Field2
}
\`\`\`

### 1.3 Repository Interface (`repository.go`)

\`\`\`go
package {entity}

import "context"

// {Entity}Repository defines the interface for {entity} persistence.
// Implementations must handle database operations and error cases.
type {Entity}Repository interface {
    // Save persists a {Entity} entity.
    // If entity exists (same ID), it should be updated.
    Save(ctx context.Context, entity *{Entity}) error

    // FindByID retrieves a {Entity} by its ID.
    // Returns ErrNotFound if entity doesn't exist.
    FindByID(ctx context.Context, id {Entity}ID) (*{Entity}, error)

    // FindAll retrieves all {Entity} entities.
    // Consider pagination for large datasets.
    FindAll(ctx context.Context) ([]*{Entity}, error)

    // Delete removes a {Entity} by ID.
    Delete(ctx context.Context, id {Entity}ID) error
}

// Domain errors
var (
    Err{Entity}NotFound = errors.New("{entity} not found")
    Err{Entity}Invalid  = errors.New("invalid {entity}")
)
\`\`\`

### 1.4 Domain Service (`service.go`) (if needed)

\`\`\`go
package {entity}

// {Entity}Service contains domain logic that doesn't belong to a single entity.
// Use sparingly - most logic should be in entities.
type {Entity}Service struct {
    // dependencies (other repositories, if needed)
}

func New{Entity}Service(/* deps */) *{Entity}Service {
    return &{Entity}Service{
        // initialize
    }
}

// Complex business operations that involve multiple entities
\`\`\`

## Step 2: Implement Repository

**Location:** `backend/internal/infrastructure/repository/mysql/`

### 2.1 Repository Implementation (`{entity}_repository.go`)

\`\`\`go
package mysql

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/xiajiayi/ai-motion/internal/domain/{entity}"
)

type MySQL{Entity}Repository struct {
    db *sql.DB
}

func NewMySQL{Entity}Repository(db *sql.DB) {entity}.{Entity}Repository {
    return &MySQL{Entity}Repository{db: db}
}

func (r *MySQL{Entity}Repository) Save(ctx context.Context, e *{entity}.{Entity}) error {
    query := \`
        INSERT INTO {entities} (id, /* fields */, created_at, updated_at)
        VALUES (?, /* placeholders */, ?, ?)
        ON DUPLICATE KEY UPDATE
            /* field = VALUES(field), */
            updated_at = VALUES(updated_at)
    \`

    _, err := r.db.ExecContext(ctx, query,
        e.ID,
        // field values
        e.CreatedAt,
        e.UpdatedAt,
    )

    if err != nil {
        return fmt.Errorf("failed to save {entity}: %w", err)
    }

    return nil
}

func (r *MySQL{Entity}Repository) FindByID(ctx context.Context, id {entity}.{Entity}ID) (*{entity}.{Entity}, error) {
    query := \`
        SELECT id, /* fields */, created_at, updated_at
        FROM {entities}
        WHERE id = ?
    \`

    var e {entity}.{Entity}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &e.ID,
        // field pointers
        &e.CreatedAt,
        &e.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, {entity}.Err{Entity}NotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find {entity}: %w", err)
    }

    return &e, nil
}

func (r *MySQL{Entity}Repository) FindAll(ctx context.Context) ([]*{entity}.{Entity}, error) {
    query := \`SELECT id, /* fields */, created_at, updated_at FROM {entities}\`

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query {entities}: %w", err)
    }
    defer rows.Close()

    var entities []*{entity}.{Entity}
    for rows.Next() {
        var e {entity}.{Entity}
        if err := rows.Scan(&e.ID, /* fields */, &e.CreatedAt, &e.UpdatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan {entity}: %w", err)
        }
        entities = append(entities, &e)
    }

    return entities, nil
}

func (r *MySQL{Entity}Repository) Delete(ctx context.Context, id {entity}.{Entity}ID) error {
    query := \`DELETE FROM {entities} WHERE id = ?\`

    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete {entity}: %w", err)
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return {entity}.Err{Entity}NotFound
    }

    return nil
}
\`\`\`

## Step 3: Create Application Layer

**Location:** `backend/internal/application/`

### 3.1 DTOs (`dto/{entity}_dto.go`)

\`\`\`go
package dto

import "time"

// {Entity}Request represents the request to create/update a {entity}.
type {Entity}Request struct {
    // Request fields (from client)
}

// {Entity}Response represents the response containing {entity} data.
type {Entity}Response struct {
    ID        string    \`json:"id"\`
    // Response fields
    CreatedAt time.Time \`json:"created_at"\`
    UpdatedAt time.Time \`json:"updated_at"\`
}
\`\`\`

### 3.2 Application Service (`service/{entity}_service.go`)

\`\`\`go
package service

import (
    "context"
    "fmt"
    "github.com/xiajiayi/ai-motion/internal/domain/{entity}"
    "github.com/xiajiayi/ai-motion/internal/application/dto"
)

type {Entity}Service struct {
    {entity}Repo {entity}.{Entity}Repository
}

func New{Entity}Service({entity}Repo {entity}.{Entity}Repository) *{Entity}Service {
    return &{Entity}Service{
        {entity}Repo: {entity}Repo,
    }
}

func (s *{Entity}Service) Create(ctx context.Context, req *dto.{Entity}Request) (*dto.{Entity}Response, error) {
    // 1. Create domain entity
    e := {entity}.New{Entity}(/* map from req */)

    // 2. Validate
    if err := e.Validate(); err != nil {
        return nil, fmt.Errorf("invalid {entity}: %w", err)
    }

    // 3. Persist
    if err := s.{entity}Repo.Save(ctx, e); err != nil {
        return nil, fmt.Errorf("failed to save {entity}: %w", err)
    }

    // 4. Return DTO
    return s.toDTO(e), nil
}

func (s *{Entity}Service) GetByID(ctx context.Context, id string) (*dto.{Entity}Response, error) {
    e, err := s.{entity}Repo.FindByID(ctx, {entity}.{Entity}ID(id))
    if err != nil {
        return nil, err
    }

    return s.toDTO(e), nil
}

func (s *{Entity}Service) List(ctx context.Context) ([]*dto.{Entity}Response, error) {
    entities, err := s.{entity}Repo.FindAll(ctx)
    if err != nil {
        return nil, err
    }

    responses := make([]*dto.{Entity}Response, len(entities))
    for i, e := range entities {
        responses[i] = s.toDTO(e)
    }

    return responses, nil
}

func (s *{Entity}Service) toDTO(e *{entity}.{Entity}) *dto.{Entity}Response {
    return &dto.{Entity}Response{
        ID:        string(e.ID),
        // Map fields
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
    }
}
\`\`\`

## Step 4: Create Interface Layer

**Location:** `backend/internal/interfaces/http/handler/`

### 4.1 Handler (`{entity}_handler.go`)

\`\`\`go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/xiajiayi/ai-motion/internal/application/service"
    "github.com/xiajiayi/ai-motion/internal/application/dto"
)

type {Entity}Handler struct {
    {entity}Service *service.{Entity}Service
}

func New{Entity}Handler({entity}Service *service.{Entity}Service) *{Entity}Handler {
    return &{Entity}Handler{
        {entity}Service: {entity}Service,
    }
}

func (h *{Entity}Handler) Create(c *gin.Context) {
    var req dto.{Entity}Request

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request",
            "details": err.Error(),
        })
        return
    }

    resp, err := h.{entity}Service.Create(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create {entity}",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *{Entity}Handler) GetByID(c *gin.Context) {
    id := c.Param("id")

    resp, err := h.{entity}Service.GetByID(c.Request.Context(), id)
    if err != nil {
        status := http.StatusInternalServerError
        if errors.Is(err, {entity}.Err{Entity}NotFound) {
            status = http.StatusNotFound
        }
        c.JSON(status, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *{Entity}Handler) List(c *gin.Context) {
    resp, err := h.{entity}Service.List(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": resp})
}
\`\`\`

## Step 5: Register Routes

**Location:** `backend/cmd/main.go`

\`\`\`go
// Initialize repository
{entity}Repo := mysql.NewMySQL{Entity}Repository(db)

// Initialize service
{entity}Service := service.New{Entity}Service({entity}Repo)

// Initialize handler
{entity}Handler := handler.New{Entity}Handler({entity}Service)

// Register routes
{entity}Group := v1.Group("/{entities}")
{
    {entity}Group.POST("/", {entity}Handler.Create)
    {entity}Group.GET("/:id", {entity}Handler.GetByID)
    {entity}Group.GET("/", {entity}Handler.List)
}
\`\`\`

## Step 6: Add Tests

See the `add-tests` skill for detailed testing patterns.

## Checklist

After implementation, verify:

- [ ] Entity is in domain layer with business logic
- [ ] Repository interface is in domain layer
- [ ] Repository implementation is in infrastructure layer
- [ ] Application service orchestrates use cases
- [ ] DTOs are used for data transfer
- [ ] Handler validates input and formats response
- [ ] Routes are registered in main.go
- [ ] Unit tests for domain logic
- [ ] Integration tests for repository
- [ ] Handler tests with mocks
- [ ] No domain layer imports infrastructure
- [ ] All errors are properly wrapped
- [ ] Documentation comments added

Remember: Follow DDD strictly. Domain layer should be pure business logic with zero infrastructure dependencies.
