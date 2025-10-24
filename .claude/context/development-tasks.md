# Common Development Tasks

## Adding a New Domain Entity

1. Define entity in `backend/internal/domain/{entity}/entity.go`
2. Define repository interface in same package
3. Implement in `backend/internal/infrastructure/repository/mysql/`
4. Create application service in `backend/internal/application/service/`
5. Create DTOs in `backend/internal/application/dto/`
6. Implement handler in `backend/internal/interfaces/http/handler/`
7. Register routes in `backend/cmd/main.go`

## Adding a New API Endpoint

1. Define DTOs in application layer
2. Implement handler in `backend/internal/interfaces/http/handler/`
3. Add route in `backend/cmd/main.go`
4. Add middleware if needed
5. Update `docs/API.md`
6. Write tests

## Creating a New Frontend Page

1. Create page component in `frontend/src/pages/`
2. Create API service in `frontend/src/services/`
3. Define TypeScript interfaces in `frontend/src/types/`
4. Add route in `frontend/src/App.tsx`
5. Create sub-components in `frontend/src/components/features/`
