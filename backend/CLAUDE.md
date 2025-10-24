# Backend CLAUDE.md

This file provides backend-specific guidance for Claude Code when working with the Go backend.

## Backend Overview

The AI-Motion backend is a Go 1.24+ application using the Gin web framework with strict Domain-Driven Design (DDD) architecture.

**Module:** `github.com/xiajiayi/ai-motion`
**Framework:** Gin
**Architecture:** DDD with 4 layers
**Current State:** Basic structure in place, implementing full DDD pattern

## Directory Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point, route registration
├── internal/
│   ├── domain/              # Domain Layer - Core business logic
│   │   ├── novel/          # Novel aggregate root
│   │   ├── character/      # Character aggregate root
│   │   ├── scene/          # Scene aggregate root
│   │   └── media/          # Media aggregate root
│   ├── application/         # Application Layer - Use cases
│   │   ├── service/        # Application services
│   │   └── dto/            # Data Transfer Objects
│   ├── infrastructure/      # Infrastructure Layer - Technical concerns
│   │   ├── repository/     # Repository implementations
│   │   │   └── mysql/      # MySQL implementations
│   │   ├── ai/             # AI service clients
│   │   │   ├── gemini/     # Gemini API client
│   │   │   └── sora/       # Sora2 API client
│   │   ├── storage/        # File storage (local/MinIO)
│   │   └── config/         # Configuration management
│   └── interfaces/          # Interface Layer - External interactions
│       ├── http/
│       │   ├── handler/    # HTTP handlers
│       │   ├── request/    # Request validation
│       │   └── response/   # Response formatting
│       └── middleware/      # Middleware (auth, logging, CORS)
├── pkg/                     # Public packages (reusable utilities)
│   ├── ai/                 # AI client interfaces
│   ├── storage/            # Storage interfaces
│   └── utils/              # Utility functions
├── storage/                 # Local file storage (gitignored)
├── config/                  # Configuration files
│   └── config.example.yaml
├── go.mod
└── go.sum
```

## Development Commands

### Running the Backend

```bash
# Development mode (from backend directory)
go run cmd/main.go

# Or from project root
cd backend && go run cmd/main.go

# With hot reload (install air first)
go install github.com/air-verse/air@latest
air

# Production build
go build -o bin/ai-motion cmd/main.go
./bin/ai-motion
```

### Dependencies

```bash
# Download dependencies
go mod download

# Add new dependency
go get github.com/some/package@latest

# Update dependencies
go get -u ./...

# Tidy up go.mod and go.sum
go mod tidy

# Verify dependencies
go mod verify
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run tests for specific package
go test ./internal/domain/novel

# Run specific test
go test ./internal/domain/novel -run TestNovelCreation

# Run tests with verbose output
go test -v ./...

# Run tests with race detector
go test -race ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code (install golangci-lint first)
golangci-lint run

# Vet code
go vet ./...

# Check for common mistakes
staticcheck ./...
```

## DDD Architecture Guidelines

### Layer Responsibilities

**1. Domain Layer (`internal/domain/`)**
- **Contains:** Entities, Value Objects, Repository Interfaces, Domain Services
- **Dependencies:** NONE - Must be pure business logic
- **Rules:**
  - No external package imports (except standard library)
  - No database, HTTP, or infrastructure code
  - All business rules live here
  - Repository interfaces defined here, implemented elsewhere

**Example Domain Entity:**
```go
// internal/domain/novel/entity.go
package novel

type Novel struct {
    ID        NovelID
    Title     string
    Author    string
    Content   string
    Chapters  []Chapter
    Status    NovelStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Business logic method
func (n *Novel) Parse() error {
    // Domain logic for parsing novel
}

func (n *Novel) ExtractCharacters() ([]CharacterID, error) {
    // Domain logic for character extraction
}
```

**Example Repository Interface:**
```go
// internal/domain/novel/repository.go
package novel

type NovelRepository interface {
    Save(ctx context.Context, novel *Novel) error
    FindByID(ctx context.Context, id NovelID) (*Novel, error)
    FindAll(ctx context.Context) ([]*Novel, error)
    Delete(ctx context.Context, id NovelID) error
}
```

**2. Application Layer (`internal/application/`)**
- **Contains:** Application Services, DTOs, Use Case Orchestration
- **Dependencies:** Domain Layer only
- **Rules:**
  - Orchestrates domain objects to fulfill use cases
  - Manages transactions
  - Converts between domain models and DTOs
  - No business logic (delegate to domain)

**Example Application Service:**
```go
// internal/application/service/novel_service.go
package service

type NovelService struct {
    novelRepo      domain.NovelRepository
    characterRepo  domain.CharacterRepository
    parserService  domain.NovelParserService
}

func NewNovelService(
    novelRepo domain.NovelRepository,
    characterRepo domain.CharacterRepository,
    parserService domain.NovelParserService,
) *NovelService {
    return &NovelService{
        novelRepo:     novelRepo,
        characterRepo: characterRepo,
        parserService: parserService,
    }
}

func (s *NovelService) UploadAndParse(ctx context.Context, req *dto.UploadNovelRequest) (*dto.NovelResponse, error) {
    // 1. Create domain entity
    novel := domain.NewNovel(req.Title, req.Author, req.Content)

    // 2. Execute domain logic
    if err := novel.Parse(); err != nil {
        return nil, fmt.Errorf("failed to parse novel: %w", err)
    }

    // 3. Extract characters
    characterIDs, err := novel.ExtractCharacters()
    if err != nil {
        return nil, fmt.Errorf("failed to extract characters: %w", err)
    }

    // 4. Persist (with transaction if needed)
    if err := s.novelRepo.Save(ctx, novel); err != nil {
        return nil, fmt.Errorf("failed to save novel: %w", err)
    }

    // 5. Return DTO
    return s.toDTO(novel, characterIDs), nil
}

func (s *NovelService) toDTO(novel *domain.Novel, charIDs []domain.CharacterID) *dto.NovelResponse {
    return &dto.NovelResponse{
        ID:           string(novel.ID),
        Title:        novel.Title,
        Author:       novel.Author,
        Status:       string(novel.Status),
        CharacterIDs: charIDs,
        CreatedAt:    novel.CreatedAt,
    }
}
```

**3. Infrastructure Layer (`internal/infrastructure/`)**
- **Contains:** Repository implementations, External service clients, Storage, Config
- **Dependencies:** Domain and Application layers
- **Rules:**
  - Implements repository interfaces from domain
  - Handles all external communications
  - Database access, API calls, file I/O
  - Configuration loading

**Example Repository Implementation:**
```go
// internal/infrastructure/repository/mysql/novel_repository.go
package mysql

type MySQLNovelRepository struct {
    db *sql.DB
}

func NewMySQLNovelRepository(db *sql.DB) domain.NovelRepository {
    return &MySQLNovelRepository{db: db}
}

func (r *MySQLNovelRepository) Save(ctx context.Context, novel *domain.Novel) error {
    query := `
        INSERT INTO novels (id, title, author, content, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            title = VALUES(title),
            content = VALUES(content),
            status = VALUES(status),
            updated_at = VALUES(updated_at)
    `

    _, err := r.db.ExecContext(ctx, query,
        novel.ID, novel.Title, novel.Author, novel.Content,
        novel.Status, novel.CreatedAt, novel.UpdatedAt,
    )

    if err != nil {
        return fmt.Errorf("failed to save novel: %w", err)
    }

    return nil
}

func (r *MySQLNovelRepository) FindByID(ctx context.Context, id domain.NovelID) (*domain.Novel, error) {
    query := `
        SELECT id, title, author, content, status, created_at, updated_at
        FROM novels WHERE id = ?
    `

    var novel domain.Novel
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &novel.ID, &novel.Title, &novel.Author, &novel.Content,
        &novel.Status, &novel.CreatedAt, &novel.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, domain.ErrNovelNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find novel: %w", err)
    }

    return &novel, nil
}
```

**4. Interface Layer (`internal/interfaces/`)**
- **Contains:** HTTP handlers, Middleware, Request/Response types
- **Dependencies:** Application layer
- **Rules:**
  - Handles HTTP requests/responses
  - Request validation
  - Error formatting
  - Authentication/Authorization

**Example Handler:**
```go
// internal/interfaces/http/handler/novel_handler.go
package handler

type NovelHandler struct {
    novelService *application.NovelService
}

func NewNovelHandler(novelService *application.NovelService) *NovelHandler {
    return &NovelHandler{novelService: novelService}
}

func (h *NovelHandler) Upload(c *gin.Context) {
    var req dto.UploadNovelRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request",
            "details": err.Error(),
        })
        return
    }

    // Validate request
    if req.Title == "" || req.Content == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Title and content are required",
        })
        return
    }

    // Call application service
    novel, err := h.novelService.UploadAndParse(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to upload novel",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": novel,
    })
}
```

## Coding Standards

### Naming Conventions

```go
// Packages: lowercase, singular
package novel

// Interfaces: noun or adjective
type NovelRepository interface {}
type Parsable interface {}

// Structs: PascalCase
type Novel struct {}
type CharacterService struct {}

// Methods: PascalCase (exported) or camelCase (private)
func (n *Novel) Parse() error {}
func (n *Novel) validateContent() error {}

// Variables: camelCase
var novelCount int
var maxCharacters = 1000

// Constants: PascalCase or UPPER_SNAKE_CASE
const MaxChapterSize = 10000
const DEFAULT_PAGE_SIZE = 20
```

### Error Handling

```go
// Always check errors
result, err := someOperation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Define domain errors
var (
    ErrNovelNotFound = errors.New("novel not found")
    ErrInvalidContent = errors.New("invalid content")
)

// Wrap errors with context
if err := repo.Save(ctx, novel); err != nil {
    return fmt.Errorf("failed to save novel %s: %w", novel.ID, err)
}

// Use errors.Is and errors.As
if errors.Is(err, domain.ErrNovelNotFound) {
    return http.StatusNotFound
}
```

### Documentation

```go
// Package documentation
// Package novel provides novel parsing and management functionality.
// It implements the Novel aggregate root in the DDD architecture.
package novel

// Type documentation
// Novel represents a novel document with chapters and metadata.
// It is an aggregate root in the domain model.
type Novel struct {
    ID     NovelID
    Title  string
}

// Function documentation
// ParseNovel parses a novel file and extracts chapters and characters.
// It returns an error if the file format is invalid or parsing fails.
//
// Example:
//   novel, err := ParseNovel(reader)
//   if err != nil {
//       log.Fatal(err)
//   }
func ParseNovel(r io.Reader) (*Novel, error) {
    // Implementation
}
```

### Context Usage

```go
// Always accept context as first parameter
func (s *NovelService) GetNovel(ctx context.Context, id string) (*Novel, error) {
    // Use context for cancellation and deadlines
    if err := ctx.Err(); err != nil {
        return nil, err
    }

    // Pass context to downstream calls
    return s.repo.FindByID(ctx, id)
}

// Set timeout for external API calls
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

result, err := externalAPI.Call(ctx, params)
```

## Testing Guidelines

### Unit Tests

```go
// internal/domain/novel/novel_test.go
package novel_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/xiajiayi/ai-motion/internal/domain/novel"
)

func TestNovel_Parse(t *testing.T) {
    tests := []struct {
        name    string
        content string
        wantErr bool
    }{
        {
            name:    "valid content",
            content: "Chapter 1\nSome content...",
            wantErr: false,
        },
        {
            name:    "empty content",
            content: "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            n := novel.NewNovel("Test", "Author", tt.content)
            err := n.Parse()

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Integration Tests

```go
// internal/infrastructure/repository/mysql/novel_repository_test.go
package mysql_test

import (
    "context"
    "testing"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("mysql", "test_dsn")
    require.NoError(t, err)

    // Run migrations
    // ...

    t.Cleanup(func() {
        db.Close()
    })

    return db
}

func TestMySQLNovelRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := mysql.NewMySQLNovelRepository(db)

    novel := domain.NewNovel("Test", "Author", "Content")
    err := repo.Save(context.Background(), novel)

    assert.NoError(t, err)
}
```

### Mock Repositories

```go
// Use interfaces for easy mocking
type MockNovelRepository struct {
    SaveFunc    func(ctx context.Context, novel *domain.Novel) error
    FindByIDFunc func(ctx context.Context, id domain.NovelID) (*domain.Novel, error)
}

func (m *MockNovelRepository) Save(ctx context.Context, novel *domain.Novel) error {
    if m.SaveFunc != nil {
        return m.SaveFunc(ctx, novel)
    }
    return nil
}

func (m *MockNovelRepository) FindByID(ctx context.Context, id domain.NovelID) (*domain.Novel, error) {
    if m.FindByIDFunc != nil {
        return m.FindByIDFunc(ctx, id)
    }
    return nil, nil
}
```

## AI Service Integration

### Gemini Client Structure

```go
// internal/infrastructure/ai/gemini/client.go
package gemini

type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func NewClient(apiKey string) *Client {
    return &Client{
        apiKey:     apiKey,
        baseURL:    "https://api.gemini.google.com/v1",
        httpClient: &http.Client{Timeout: 60 * time.Second},
    }
}

func (c *Client) TextToImage(ctx context.Context, prompt string) (string, error) {
    // Implementation
}

func (c *Client) ImageToImage(ctx context.Context, imageURL, prompt string) (string, error) {
    // Implementation
}
```

### Sora2 Client Structure

```go
// internal/infrastructure/ai/sora/client.go
package sora

type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func (c *Client) ImageToVideo(ctx context.Context, imageURL string) (string, error) {
    // Implementation
}
```

## Database Patterns

### Connection Management

```go
// internal/infrastructure/database/mysql.go
func NewMySQLConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
```

### Transaction Management

```go
func (s *NovelService) CreateWithCharacters(ctx context.Context, req *dto.CreateNovelRequest) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Save novel
    if err := s.novelRepo.SaveTx(ctx, tx, novel); err != nil {
        return err
    }

    // Save characters
    for _, char := range characters {
        if err := s.characterRepo.SaveTx(ctx, tx, char); err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

## Performance Best Practices

### Concurrency

```go
// Use goroutines for concurrent operations
func (s *GenerationService) GenerateMultipleScenes(ctx context.Context, sceneIDs []string) error {
    errCh := make(chan error, len(sceneIDs))

    for _, id := range sceneIDs {
        go func(sceneID string) {
            err := s.GenerateScene(ctx, sceneID)
            errCh <- err
        }(id)
    }

    // Collect errors
    for i := 0; i < len(sceneIDs); i++ {
        if err := <-errCh; err != nil {
            return err
        }
    }

    return nil
}

// Use worker pool for rate limiting
func (s *Service) ProcessWithRateLimit(ctx context.Context, items []Item) error {
    sem := make(chan struct{}, 5) // Max 5 concurrent

    for _, item := range items {
        sem <- struct{}{} // Acquire

        go func(i Item) {
            defer func() { <-sem }() // Release
            s.process(ctx, i)
        }(item)
    }

    // Wait for all to complete
    for i := 0; i < cap(sem); i++ {
        sem <- struct{}{}
    }

    return nil
}
```

### Caching

```go
// Simple in-memory cache
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

## Debugging Tips

```go
// Use structured logging
import "log/slog"

slog.Info("Processing novel",
    "novel_id", novel.ID,
    "title", novel.Title,
    "chapters", len(novel.Chapters))

slog.Error("Failed to save",
    "error", err,
    "novel_id", novel.ID)

// Use Delve debugger
// Install: go install github.com/go-delve/delve/cmd/dlv@latest
// Run: dlv debug cmd/main.go
// Set breakpoint: break main.main
// Continue: continue
// Step: next, step
// Inspect: print variableName
```

## Common Pitfalls to Avoid

1. **Don't** put business logic in handlers - Use domain/application layers
2. **Don't** import infrastructure packages in domain layer
3. **Don't** ignore errors - Always check and handle them
4. **Don't** use `panic` for error handling - Return errors
5. **Don't** share state between requests - Keep handlers stateless
6. **Don't** forget to close resources - Use `defer`
7. **Don't** block forever - Always use context with timeout
8. **Don't** commit API keys - Use environment variables

## Useful Go Tools

```bash
# Code generation
go generate ./...

# Dependency visualization
go mod graph | grep "your-package"

# Find unused dependencies
go mod tidy

# Security check
gosec ./...

# Benchmarking
go test -bench=. -benchmem ./...
```

## Quick Reference

**Current Routes:** See [cmd/main.go:22](cmd/main.go#L22)
**Database Schema:** See [docs/ARCHITECTURE.md:554](docs/ARCHITECTURE.md#L554)
**API Docs:** See [docs/API.md](docs/API.md)
**Main Architecture:** See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
