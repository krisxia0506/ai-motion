---
name: add-tests
description: Generate comprehensive unit, integration, and E2E tests for backend and frontend
version: 1.0.0
author: AI-Motion Team
---

# Add Tests Skill

You are an expert in writing comprehensive tests for the AI-Motion project.

## Testing Guidelines

### Backend (Go) Testing

#### Unit Tests for Domain Logic

**Location:** `internal/domain/{entity}/{entity}_test.go`

**Pattern:**
\`\`\`go
package novel_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/xiajiayi/ai-motion/internal/domain/novel"
)

func TestNovel_MethodName(t *testing.T) {
    // Table-driven tests
    tests := []struct {
        name    string
        input   inputType
        want    expectedType
        wantErr bool
    }{
        {
            name:    "valid case",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "error case",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            entity := setupEntity()

            // Execute
            got, err := entity.Method(tt.input)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
\`\`\`

#### Integration Tests for Repositories

**Location:** `internal/infrastructure/repository/mysql/{entity}_repository_test.go`

**Pattern:**
\`\`\`go
package mysql_test

import (
    "context"
    "testing"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
    // Setup test database
    db, err := sql.Open("mysql", "test_dsn")
    require.NoError(t, err)

    // Run migrations or create tables

    // Cleanup after test
    t.Cleanup(func() {
        // Clean test data
        db.Exec("TRUNCATE TABLE novels")
        db.Close()
    })

    return db
}

func TestMySQLNovelRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := mysql.NewMySQLNovelRepository(db)

    novel := domain.NewNovel("Title", "Author", "Content")

    err := repo.Save(context.Background(), novel)
    assert.NoError(t, err)

    // Verify saved
    found, err := repo.FindByID(context.Background(), novel.ID)
    require.NoError(t, err)
    assert.Equal(t, novel.Title, found.Title)
}
\`\`\`

#### Handler Tests with Mocks

**Location:** `internal/interfaces/http/handler/{entity}_handler_test.go`

**Pattern:**
\`\`\`go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockNovelService struct {
    mock.Mock
}

func (m *MockNovelService) Upload(ctx context.Context, req *dto.UploadRequest) (*dto.NovelResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.NovelResponse), args.Error(1)
}

func TestNovelHandler_Upload(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockService := new(MockNovelService)
    handler := NewNovelHandler(mockService)

    t.Run("success", func(t *testing.T) {
        // Setup mock
        expectedResponse := &dto.NovelResponse{ID: "123", Title: "Test"}
        mockService.On("Upload", mock.Anything, mock.Anything).Return(expectedResponse, nil)

        // Create request
        reqBody := map[string]string{"title": "Test", "content": "Content"}
        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequest(http.MethodPost, "/api/v1/novel/upload", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()

        // Execute
        router := gin.New()
        router.POST("/api/v1/novel/upload", handler.Upload)
        router.ServeHTTP(w, req)

        // Assert
        assert.Equal(t, http.StatusOK, w.Code)
        mockService.AssertExpectations(t)
    })
}
\`\`\`

### Frontend (React) Testing

#### Component Tests

**Location:** `src/components/{feature}/{Component}.test.tsx`

**Pattern:**
\`\`\`typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { NovelCard } from './NovelCard';

describe('NovelCard', () => {
    const mockNovel = {
        id: '1',
        title: 'Test Novel',
        author: 'Test Author',
        status: 'completed' as const,
    };

    it('renders novel information', () => {
        render(<NovelCard novel={mockNovel} onSelect={() => {}} />);

        expect(screen.getByText('Test Novel')).toBeInTheDocument();
        expect(screen.getByText(/Test Author/i)).toBeInTheDocument();
    });

    it('calls onSelect when clicked', () => {
        const handleSelect = jest.fn();
        render(<NovelCard novel={mockNovel} onSelect={handleSelect} />);

        fireEvent.click(screen.getByText('Test Novel'));

        expect(handleSelect).toHaveBeenCalledWith('1');
    });

    it('shows loading state', () => {
        render(<NovelCard novel={mockNovel} loading={true} onSelect={() => {}} />);

        expect(screen.getByText(/loading/i)).toBeInTheDocument();
    });
});
\`\`\`

#### Hook Tests

**Location:** `src/hooks/{hookName}.test.ts`

**Pattern:**
\`\`\`typescript
import { renderHook, waitFor } from '@testing-library/react';
import { useNovel } from './useNovel';
import * as novelApi from '../services/novelApi';

jest.mock('../services/novelApi');

describe('useNovel', () => {
    it('fetches novel data', async () => {
        const mockNovel = { id: '1', title: 'Test' };
        (novelApi.getNovel as jest.Mock).mockResolvedValue(mockNovel);

        const { result } = renderHook(() => useNovel('1'));

        expect(result.current.loading).toBe(true);

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.novel).toEqual(mockNovel);
        expect(result.current.error).toBeNull();
    });

    it('handles errors', async () => {
        const error = new Error('Fetch failed');
        (novelApi.getNovel as jest.Mock).mockRejectedValue(error);

        const { result } = renderHook(() => useNovel('1'));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toEqual(error);
        expect(result.current.novel).toBeNull();
    });
});
\`\`\`

#### Integration Tests

**Location:** `src/__tests__/integration/{feature}.test.tsx`

**Pattern:**
\`\`\`typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { NovelListPage } from '../pages/NovelListPage';
import * as novelApi from '../services/novelApi';

jest.mock('../services/novelApi');

describe('NovelListPage Integration', () => {
    it('displays novels and handles selection', async () => {
        const mockNovels = [
            { id: '1', title: 'Novel 1' },
            { id: '2', title: 'Novel 2' },
        ];
        (novelApi.listNovels as jest.Mock).mockResolvedValue(mockNovels);

        render(
            <BrowserRouter>
                <NovelListPage />
            </BrowserRouter>
        );

        // Wait for loading to finish
        await waitFor(() => {
            expect(screen.getByText('Novel 1')).toBeInTheDocument();
        });

        // Click novel
        fireEvent.click(screen.getByText('Novel 1'));

        // Verify navigation or state change
        // ...
    });
});
\`\`\`

## Test Coverage Requirements

### Backend
- **Domain layer:** 80%+ coverage (critical business logic)
- **Application layer:** 70%+ coverage
- **Handlers:** 60%+ coverage (happy path + error cases)
- **Infrastructure:** Integration tests for repository CRUD

### Frontend
- **Custom hooks:** 80%+ coverage
- **Feature components:** 70%+ coverage
- **Common components:** 60%+ coverage
- **Integration tests:** Critical user flows

## Commands to Run Tests

### Backend
\`\`\`bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test ./internal/domain/novel -run TestNovel_Parse

# Run with race detector
go test -race ./...
\`\`\`

### Frontend
\`\`\`bash
# Run all tests
npm test

# Run in watch mode
npm test -- --watch

# Run with coverage
npm test -- --coverage

# Run specific test file
npm test NovelCard.test.tsx
\`\`\`

## Test Writing Checklist

When writing tests:

- [ ] Test the happy path (success case)
- [ ] Test error cases (failure scenarios)
- [ ] Test edge cases (empty input, null, boundary values)
- [ ] Test concurrent operations (if applicable)
- [ ] Mock external dependencies (APIs, databases)
- [ ] Use descriptive test names
- [ ] Cleanup resources in test teardown
- [ ] Don't test implementation details, test behavior
- [ ] Keep tests isolated (no shared state)
- [ ] Tests should be fast and reliable

## When to Write Tests

**Always write tests for:**
- Domain entities business logic
- Complex algorithms or calculations
- Critical user flows
- Bug fixes (regression tests)

**Tests are optional for:**
- Simple CRUD operations with no logic
- Trivial getters/setters
- UI layout components with no logic

Remember: Good tests document behavior and prevent regressions. Write tests that would help future developers understand what the code does.