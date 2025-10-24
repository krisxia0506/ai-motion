# Claude Skills for AI-Motion

This directory contains reusable Claude skills for common development tasks in the AI-Motion project.

## What are Claude Skills?

Skills are specialized prompt templates that give Claude specific expertise and instructions for particular tasks. They help ensure consistent, high-quality code and follow project conventions.

## Skill Structure

Each skill is organized as a folder containing:
- `SKILL.md` - The main skill prompt with instructions, examples, and checklists
- (Optional) Additional resources like templates, examples, or configuration files

```
.claude/skills/
├── review-code/
│   └── SKILL.md
├── add-tests/
│   └── SKILL.md
└── implement-ddd-entity/
    └── SKILL.md
```

## Available Skills

### 🔍 review-code
**Location:** [.claude/skills/review-code/](.claude/skills/review-code/)

Reviews code changes for:
- DDD architecture compliance
- Code quality and best practices
- Security issues
- Performance concerns
- Testing coverage
- Project-specific patterns

**When to use:**
- Before committing changes
- During pull request reviews
- When you want feedback on new code

**Example:**
```
@review-code

Please review these changes I made to the Novel entity and handler.
[paste code or reference files]
```

---

### ✅ add-tests
**Location:** [.claude/skills/add-tests/](.claude/skills/add-tests/)

Generates comprehensive tests for:
- Backend: Unit tests, integration tests, handler tests
- Frontend: Component tests, hook tests, integration tests
- Follows project testing patterns
- Includes mocks and test utilities

**When to use:**
- After implementing new features
- When test coverage is low
- To see testing best practices

**Example:**
```
@add-tests

I just implemented the Novel entity and service. Please generate complete tests for:
- internal/domain/novel/entity.go
- internal/application/service/novel_service.go
```

---

### 🏗️ implement-ddd-entity
**Location:** [.claude/skills/implement-ddd-entity/](.claude/skills/implement-ddd-entity/)

Guides complete DDD entity implementation:
- Domain layer (entity, value objects, repository interface)
- Infrastructure layer (repository implementation)
- Application layer (DTOs, service)
- Interface layer (handlers)
- Route registration
- Testing setup

**When to use:**
- Creating new domain entities
- When unsure about DDD structure
- To ensure architectural compliance

**Example:**
```
@implement-ddd-entity

I need to implement a new "Project" entity that:
- Has title, description, owner
- Can be published/unpublished
- Belongs to a user
- Has multiple scenes

Please guide me through the complete implementation.
```

## How to Use Skills

### Method 1: Direct Reference (Recommended)

In your conversation with Claude, use the `@` symbol followed by the skill name:

```
@review-code

[your request]
```

Claude will automatically load the skill context and apply its specialized knowledge.

### Method 2: Explicit Skill Invocation

You can also explicitly tell Claude to use a skill:

```
Please use the "review-code" skill to check my implementation of the character consistency feature.
```

### Method 3: Combine Multiple Skills

You can combine skills for complex tasks:

```
@implement-ddd-entity

After implementation, please also use @add-tests to generate comprehensive tests.
```

## Skill Development Guidelines

### Creating New Skills

When creating new skills, follow this structure:

```markdown
# {Skill Name} Skill

You are an expert in {area} for the AI-Motion project.

## {Section 1}
[Detailed instructions]

## {Section 2}
[Examples and patterns]

## Checklist
- [ ] Item 1
- [ ] Item 2

## Examples

**❌ Bad:**
```code
// Bad example
```

**✅ Good:**
```code
// Good example
```
```

### Best Practices for Skills

1. **Be Specific**: Include project-specific patterns and conventions
2. **Provide Examples**: Show both good and bad examples
3. **Include Checklists**: Help ensure nothing is missed
4. **Reference Docs**: Link to relevant project documentation
5. **Keep Updated**: Update skills as patterns evolve

## Proposed Additional Skills

Consider creating these skills in the future:

- **`debug-issue`**: Systematic debugging approach for common issues
- **`add-api-endpoint`**: Complete guide for adding new API endpoints
- **`optimize-performance`**: Performance analysis and optimization
- **`migrate-database`**: Database migration patterns
- **`integrate-ai-service`**: AI service integration guide
- **`setup-frontend-component`**: React component setup with all patterns
- **`add-documentation`**: Generate comprehensive documentation

## Tips for Effective Skill Usage

1. **Be Clear**: Provide context and specific requirements
2. **Share Code**: Include relevant code snippets or file paths
3. **Ask Questions**: If the output isn't what you need, ask for refinements
4. **Iterate**: Skills provide guidelines, but you can customize the output
5. **Combine with Context**: Use skills alongside file context for best results

## Skill Metadata

Each skill should follow these conventions:

- **Folder name**: `{skill-name}` (kebab-case)
- **Main file**: `SKILL.md` (must be named exactly this)
- **Title**: `# {Skill Name} Skill` (Title Case)
- **First Line**: Define expertise: "You are an expert in..."
- **Structure**: Organized sections with clear headings
- **Examples**: Include code examples with syntax highlighting
- **Checklists**: For verification steps

## Contributing

To add a new skill:

1. Create `{skill-name}/` folder in `.claude/skills/`
2. Create `SKILL.md` inside the folder
3. Follow the structure outlined above
4. Include project-specific patterns
5. Add entry to this README
6. Test the skill with Claude
7. Update project [CLAUDE.md](../../../CLAUDE.md) if relevant

## Need Help?

If you're unsure which skill to use or how to use it:

```
I want to [describe your task]. Which skill should I use and how?
```

Claude will recommend the appropriate skill and show you how to use it effectively.
