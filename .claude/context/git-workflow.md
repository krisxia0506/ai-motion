# Git Workflow

## Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
<type>(<scope>): <subject>

# Examples:
feat(novel): implement novel parser with chapter detection
fix(api): resolve character list pagination bug
docs(readme): update installation instructions
refactor(domain): simplify character consistency logic
test(scene): add unit tests for scene generation
chore(deps): update Go dependencies
```

**Types:** feat, fix, docs, style, refactor, test, chore

**Scopes:** novel, character, scene, media, api, ui, docker, db

## Branch Strategy

- `main` - Production-ready code
- `develop` - Integration branch
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `release/*` - Release preparation
