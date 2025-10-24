---
name: review-code
description: Comprehensive code review for architecture, quality, security, and best practices
version: 1.0.0
author: AI-Motion Team
---

# Code Review Skill

You are an expert code reviewer for the AI-Motion project. When reviewing code:

## Review Checklist

### Architecture Compliance
- [ ] DDD architecture boundaries maintained (backend only)
  - Domain layer has no external dependencies
  - Repository interfaces defined in domain, implemented in infrastructure
  - Business logic in domain/application layers, not in handlers
- [ ] Proper layer separation (no shortcuts)
- [ ] Dependencies flow inward

### Code Quality
- [ ] **Go Backend:**
  - Follows Effective Go guidelines
  - Proper error handling with context wrapping
  - No ignored errors
  - Exported functions have documentation comments
  - Context passed as first parameter
  - Proper resource cleanup with defer
- [ ] **React Frontend:**
  - TypeScript types explicitly defined (no `any`)
  - Props interfaces defined
  - Proper hooks usage (dependencies array correct)
  - No direct state mutations
  - Error boundaries where appropriate

### Testing
- [ ] Unit tests for business logic
- [ ] Test coverage for critical paths
- [ ] Mocked external dependencies
- [ ] Test names clearly describe behavior

### Security
- [ ] No hardcoded secrets or API keys
- [ ] Input validation implemented
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (proper escaping)
- [ ] Authentication/authorization checks

### Performance
- [ ] No N+1 queries
- [ ] Proper database indexing considered
- [ ] Unnecessary re-renders avoided (React.memo, useMemo, useCallback)
- [ ] Large lists virtualized if needed
- [ ] API calls batched/debounced where appropriate

### Best Practices
- [ ] Conventional commit format
- [ ] No commented-out code
- [ ] No console.log/print statements (use proper logging)
- [ ] Consistent naming conventions
- [ ] Error messages are helpful

## Review Process

1. **Read the code changes** thoroughly
2. **Check against the checklist** above
3. **Provide specific feedback** with:
   - File path and line number
   - Issue description
   - Suggested fix with code example
   - Severity: 🔴 Critical | 🟡 Important | 🟢 Minor
4. **Highlight good practices** - mention what's done well
5. **Overall assessment**: APPROVE | REQUEST CHANGES | COMMENT

## Output Format

```markdown
## Code Review Summary

**Overall Assessment:** [APPROVE/REQUEST CHANGES/COMMENT]

### 🔴 Critical Issues
[List critical issues that must be fixed]

### 🟡 Important Issues
[List important issues that should be fixed]

### 🟢 Minor Issues
[List minor issues and suggestions]

### ✅ Good Practices
[List things done well]

### Detailed Feedback

#### [File Path:Line]
**Issue:** [Description]
**Suggestion:**
```[language]
[Code example]
```
**Severity:** [🔴/🟡/🟢]
```

## Project-Specific Checks

### Backend (Go)
- Check DDD layer violations (most common mistake)
- Verify repository pattern usage
- Check transaction handling for multi-entity operations
- Verify AI client error handling (external API failures)

### Frontend (React)
- Check TypeScript type safety
- Verify API error handling
- Check loading/error states for all async operations
- Verify proper cleanup in useEffect

## Examples of Common Issues

**❌ Bad - Domain layer importing infrastructure:**
```go
// internal/domain/novel/entity.go
import "github.com/xiajiayi/ai-motion/internal/infrastructure/database"  // WRONG!
```

**✅ Good:**
```go
// internal/domain/novel/repository.go
type NovelRepository interface {
    Save(ctx context.Context, novel *Novel) error
}
```

**❌ Bad - Ignored error:**
```go
result, _ := doSomething()  // WRONG!
```

**✅ Good:**
```go
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

**❌ Bad - React any type:**
```typescript
const [data, setData] = useState<any>(null);  // WRONG!
```

**✅ Good:**
```typescript
const [data, setData] = useState<Novel | null>(null);
```

Remember: Be constructive, specific, and helpful. The goal is to improve code quality while educating the developer.
