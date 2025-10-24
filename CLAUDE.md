# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI-Motion is an intelligent anime generation system that converts text novels into animated content with character consistency, AI voiceovers, and visual storytelling.

**Architecture:** Monorepo with Go backend (DDD) and React frontend
**Status:** v0.1.0-alpha (in development)

**Tech Stack:**
- Backend: Go 1.24+ with Gin framework, DDD architecture
- Frontend: React 19 + TypeScript with Vite 7
- AI Services: Gemini 2.5 Flash Image, Sora2
- Database: MySQL 8.0+
- DevOps: Docker + Docker Compose

## Repository Structure

```
ai-motion/
├── backend/           # Go backend with DDD architecture
│   └── CLAUDE.md     # Backend-specific guidance
├── frontend/          # React + TypeScript frontend
│   └── CLAUDE.md     # Frontend-specific guidance
├── docs/             # Comprehensive documentation
├── docker/           # Docker configurations
└── scripts/          # Setup and utility scripts
```

**For detailed guidance:**
- Working on backend? See [backend/CLAUDE.md](backend/CLAUDE.md)
- Working on frontend? See [frontend/CLAUDE.md](frontend/CLAUDE.md)

## Core Architecture Principles

### DDD (Domain-Driven Design) - Backend

The backend strictly follows DDD layering:

1. **Domain Layer** - Core business logic (Novel, Character, Scene, Media)
2. **Application Layer** - Use case orchestration, DTOs
3. **Infrastructure Layer** - External services (DB, AI APIs, storage)
4. **Interface Layer** - HTTP handlers, middleware

**Dependency Rule:** Always flow inward. Domain has zero external dependencies.

### Key Business Domains

**Novel Domain:** Upload, parse novels, extract metadata
**Character Domain:** Character extraction, consistency management, reference images
**Scene Domain:** Scene division, dialogue extraction, prompt generation
**Media Domain:** Image/video generation, storage management

### Critical Feature: Character Consistency

Character consistency ensures the same character maintains identical visual appearance across all generated scenes. This is achieved through reference image generation and image-to-image transformation.

**For detailed design and implementation, see [CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md)**

## Configuration

### Environment Variables

Root `.env` file (copy from `.env.example`):

Frontend `.env` (`frontend/.env`):


## Testing Strategy

**Backend:**
- Unit tests for domain logic
- Integration tests for repositories (test DB)
- Handler tests with mocked services

**Frontend:**
- Component tests with React Testing Library
- Integration tests for user flows
- E2E tests with Docker Compose environment

## AI Service Integration

### Gemini 2.5 Flash Image
- **Purpose:** Text-to-image, image-to-image
- **Use Cases:** Character references, scene generation with consistency
- **Location:** `backend/internal/infrastructure/ai/gemini/`

### Sora2
- **Purpose:** Text-to-video, image-to-video
- **Use Cases:** Scene animation, dynamic content generation
- **Location:** `backend/internal/infrastructure/ai/sora/`

**Pattern:** Define service interfaces in domain layer, implement in infrastructure layer.

## Performance Considerations

- Use Go goroutines for concurrent AI API calls
- Implement caching for character reference images
- Optimize database queries with proper indexing
- Frontend: Code splitting, lazy loading, React.memo
- Consider Redis for caching in production

## Documentation

Essential docs in `docs/` directory:

- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Complete DDD architecture, entity definitions, AI integration patterns
- **[CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md)** - Character consistency design and implementation
- **[DEVELOPMENT.md](docs/DEVELOPMENT.md)** - Development setup, coding standards, testing
- **[API.md](docs/API.md)** - API endpoint specifications
- **[QUICKSTART.md](QUICKSTART.md)** - Docker and local setup guide
- **[README.md](README.md)** - Project overview, features, roadmap

## Important Notes

- **Language:** Documentation is in Chinese, but code uses English
- **Database:** Optional in current phase; service starts without MySQL
- **Character Consistency:** Core feature - see [CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md) for design details
- **DDD Strictness:** Maintain clean architecture boundaries
- **Testing:** Write tests for all new domain logic
- **API Keys:** Never commit API keys; use environment variables

---

## Additional Context (On-Demand)

The following context is available when needed. Import specific files by mentioning them:

- **Commands & Setup:** `@.claude/context/commands.md` - Build, run, test, and docker commands
- **Git Workflow:** `@.claude/context/git-workflow.md` - Commit conventions and branching strategy
- **Development Tasks:** `@.claude/context/development-tasks.md` - Step-by-step guides for common tasks
- **Troubleshooting:** `@.claude/context/troubleshooting.md` - Common issues and solutions
- **Project Status:** `@.claude/context/project-status.md` - Current development status and roadmap
