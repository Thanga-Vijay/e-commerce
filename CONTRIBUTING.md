# Contributing Guidelines

Thank you for contributing to the E-Commerce Platform! This document provides guidelines and best practices for contributing to this project.

## Table of Contents
1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Coding Standards](#coding-standards)
5. [Testing Guidelines](#testing-guidelines)
6. [Commit Messages](#commit-messages)
7. [Pull Request Process](#pull-request-process)
8. [Code Review Guidelines](#code-review-guidelines)

---

## Code of Conduct

### Our Pledge
- Be respectful and inclusive
- Accept constructive criticism
- Focus on what's best for the project
- Show empathy towards other contributors

### Unacceptable Behavior
- Harassment or discrimination
- Trolling or insulting comments
- Publishing others' private information
- Unprofessional conduct

---

## Getting Started

### Prerequisites
1. Read [GETTING_STARTED.md](docs/GETTING_STARTED.md)
2. Set up development environment
3. Familiarize yourself with the architecture
4. Join team communication channels

### First-Time Contributors
1. Look for issues tagged `good-first-issue`
2. Ask questions in team channels
3. Start with documentation improvements
4. Gradually move to code contributions

---

## Development Workflow

### 1. Choose a Task
- Check the roadmap in [ROADMAP.md](docs/ROADMAP.md)
- Pick an issue from GitHub Issues
- Discuss with team lead if unsure

### 2. Create a Branch

```bash
# Update main branch
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/issue-123-user-authentication

# Branch naming convention:
# feature/issue-{number}-{short-description}
# bugfix/issue-{number}-{short-description}
# hotfix/issue-{number}-{short-description}
# docs/{short-description}
```

### 3. Make Changes
- Write code following coding standards
- Add/update tests
- Update documentation
- Test locally

### 4. Commit Changes
- Follow commit message guidelines
- Make atomic commits (one logical change per commit)
- Commit frequently

### 5. Push and Create PR
```bash
# Push branch
git push origin feature/issue-123-user-authentication

# Create Pull Request on GitHub
# Fill in PR template
# Link related issues
```

### 6. Address Review Comments
- Respond to all review comments
- Make requested changes
- Push updates to same branch

### 7. Merge
- PR approved by 2 reviewers
- All CI checks pass
- Squash and merge to main

---

## Coding Standards

### Go (Backend Services)

#### Project Structure
```
service-name/
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── repository/      # Database layer
│   ├── service/         # Business logic
│   ├── middleware/      # Middleware
│   └── utils/           # Utilities
├── pkg/                 # Public packages
├── migrations/          # Database migrations
├── tests/               # Integration tests
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

#### Naming Conventions
```go
// Packages: lowercase, no underscores
package userservice

// Interfaces: noun or adjective
type UserRepository interface {}
type Readable interface {}

// Structs: PascalCase
type UserService struct {}

// Functions/Methods: PascalCase (exported), camelCase (private)
func CreateUser() {}
func (s *UserService) createUser() {}

// Variables: camelCase
var userCount int
var maxRetries = 3

// Constants: PascalCase or SCREAMING_SNAKE_CASE
const MaxConnections = 100
const DEFAULT_TIMEOUT = 30
```

#### Code Style
```go
// Use gofmt and goimports
go fmt ./...
goimports -w .

// Error handling
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Context usage
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // Always pass context
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// Interface usage (small interfaces)
type UserReader interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

// Dependency injection
type UserService struct {
    repo   UserRepository
    cache  CacheClient
    logger *zap.Logger
}

func NewUserService(repo UserRepository, cache CacheClient, logger *zap.Logger) *UserService {
    return &UserService{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}
```

#### Documentation
```go
// Package documentation
// Package userservice provides user management functionality.
package userservice

// Exported function/type documentation
// CreateUser creates a new user account with the provided details.
// It validates the input, checks for duplicate emails, and stores the user.
// Returns the created user or an error if validation fails.
func CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // implementation
}
```

### TypeScript/React (Frontend)

#### Project Structure
```
frontend/
├── public/
├── src/
│   ├── assets/          # Images, fonts, etc.
│   ├── components/      # Reusable components
│   │   ├── common/
│   │   ├── layout/
│   │   └── product/
│   ├── features/        # Feature-based modules
│   │   ├── auth/
│   │   ├── products/
│   │   └── cart/
│   ├── hooks/           # Custom hooks
│   ├── pages/           # Page components
│   ├── services/        # API services
│   ├── store/           # Redux store
│   │   ├── slices/
│   │   └── store.ts
│   ├── types/           # TypeScript types
│   ├── utils/           # Utilities
│   ├── App.tsx
│   └── main.tsx
├── .eslintrc.js
├── .prettierrc
├── tsconfig.json
├── vite.config.ts
└── package.json
```

#### Naming Conventions
```typescript
// Components: PascalCase
export const UserProfile: React.FC = () => {}

// Custom hooks: camelCase starting with 'use'
export const useAuth = () => {}

// Utilities: camelCase
export const formatDate = (date: Date) => {}

// Types/Interfaces: PascalCase
export interface User {}
export type UserRole = 'admin' | 'customer'

// Constants: SCREAMING_SNAKE_CASE
export const API_BASE_URL = 'http://localhost:8080'
```

#### Component Style
```typescript
// Functional components with TypeScript
import React from 'react'

interface UserCardProps {
  user: User
  onEdit: (id: string) => void
}

export const UserCard: React.FC<UserCardProps> = ({ user, onEdit }) => {
  const handleEdit = () => {
    onEdit(user.id)
  }

  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <button onClick={handleEdit}>Edit</button>
    </div>
  )
}

// Custom hooks
import { useState, useEffect } from 'react'

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Fetch user
  }, [])

  return { user, loading }
}

// Redux Toolkit slice
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
}

const authSlice = createSlice({
  name: 'auth',
  initialState: { user: null, token: null, loading: false } as AuthState,
  reducers: {
    setUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload
    },
  },
})
```

### SQL (Database)

#### Migration Files
```sql
-- migrations/000001_create_users_table.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create table with proper naming
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(20) DEFAULT 'customer' CHECK (role IN ('customer', 'admin')),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE users IS 'User accounts for authentication';
COMMENT ON COLUMN users.role IS 'User role: customer or admin';
```

### YAML (Kubernetes)

```yaml
# kubernetes/deployments/auth-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: ecommerce
  labels:
    app: auth-service
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
        version: v1
    spec:
      containers:
      - name: auth-service
        image: auth-service:latest
        ports:
        - containerPort: 8081
          name: http
        env:
        - name: PORT
          value: "8081"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: auth-config
              key: db_host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: auth-secrets
              key: db_password
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
```

---

## Testing Guidelines

### Unit Tests (Go)

```go
// handlers/user_handler_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

// Test function
func TestCreateUser_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    user := &User{Email: "test@example.com"}
    
    mockRepo.On("Create", user).Return(nil)
    
    // Act
    err := service.CreateUser(user)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    user := &User{Email: "existing@example.com"}
    
    mockRepo.On("Create", user).Return(ErrDuplicateEmail)
    
    // Act
    err := service.CreateUser(user)
    
    // Assert
    assert.Error(t, err)
    assert.Equal(t, ErrDuplicateEmail, err)
}
```

### Component Tests (React)

```typescript
// components/UserCard.test.tsx
import { render, screen, fireEvent } from '@testing-library/react'
import { UserCard } from './UserCard'

describe('UserCard', () => {
  const mockUser = {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
  }

  it('renders user information', () => {
    render(<UserCard user={mockUser} onEdit={() => {}} />)
    
    expect(screen.getByText('John Doe')).toBeInTheDocument()
    expect(screen.getByText('john@example.com')).toBeInTheDocument()
  })

  it('calls onEdit when edit button is clicked', () => {
    const mockOnEdit = jest.fn()
    render(<UserCard user={mockUser} onEdit={mockOnEdit} />)
    
    fireEvent.click(screen.getByText('Edit'))
    
    expect(mockOnEdit).toHaveBeenCalledWith('1')
  })
})
```

### Test Coverage Requirements
- **Minimum Coverage:** 80%
- **Critical Paths:** 100% (payment, auth)
- **Run Before Commit:** Always

```bash
# Go
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# React
npm test -- --coverage
```

---

## Commit Messages

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements

### Examples

```
feat(auth): implement JWT authentication

- Add JWT token generation
- Add token validation middleware
- Add refresh token mechanism

Closes #123
```

```
fix(cart): resolve race condition in stock reservation

The cart service was not properly handling concurrent stock
reservations, leading to overselling. This fix adds proper
locking mechanism using Redis distributed locks.

Fixes #456
```

```
docs(api): update product API documentation

- Add examples for filtering
- Document pagination parameters
- Add error response codes
```

### Rules
1. Use imperative mood ("add" not "added")
2. Don't capitalize first letter
3. No period at the end of subject
4. Subject max 50 characters
5. Body max 72 characters per line
6. Reference issues/PRs in footer

---

## Pull Request Process

### Before Creating PR

- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Commit messages follow convention
- [ ] Branch is up to date with main

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## How Has This Been Tested?
Describe the tests you ran

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass
- [ ] No new warnings

## Related Issues
Closes #123
Related to #456

## Screenshots (if applicable)
```

### PR Review Process

1. **Automated Checks** (must pass)
   - Build successful
   - Tests pass
   - Linting passes
   - Coverage meets threshold

2. **Code Review** (2 approvals required)
   - Code quality
   - Architecture alignment
   - Test coverage
   - Documentation

3. **Manual Testing**
   - Reviewer tests functionality
   - Edge cases verified

4. **Approval & Merge**
   - Squash and merge to main
   - Delete branch after merge

---

## Code Review Guidelines

### As a Reviewer

#### What to Look For

**Correctness**
- Does the code do what it's supposed to?
- Are there any bugs or edge cases not handled?
- Is error handling adequate?

**Architecture**
- Does it follow project architecture?
- Is it in the right layer (handler/service/repository)?
- Does it maintain separation of concerns?

**Code Quality**
- Is it readable and maintainable?
- Are variable/function names clear?
- Is it DRY (Don't Repeat Yourself)?
- Are there comments for complex logic?

**Testing**
- Are there adequate tests?
- Do tests cover edge cases?
- Is test coverage maintained/improved?

**Performance**
- Are there any performance issues?
- Is caching used appropriately?
- Are database queries optimized?

**Security**
- Are inputs validated?
- Is data sanitized?
- Are secrets handled properly?
- Is authentication/authorization correct?

#### How to Give Feedback

**Be Constructive**
```
❌ "This code is bad"
✅ "Consider extracting this logic into a separate function for better readability"

❌ "Wrong approach"
✅ "Have you considered using X pattern here? It might be more suitable because..."
```

**Ask Questions**
```
"Could you explain why you chose this approach?"
"What's the expected behavior when X happens?"
"Have you considered edge case Y?"
```

**Suggest Alternatives**
```
"This works, but you might want to consider using the repository pattern here"
"Instead of duplicating this logic, could we create a shared utility function?"
```

**Praise Good Work**
```
"Nice refactoring!"
"Great test coverage!"
"I like how you handled this edge case"
```

### As a PR Author

#### Respond to All Comments
- Address every review comment
- Explain your reasoning if you disagree
- Ask for clarification if needed

#### Make Requested Changes
- Fix issues promptly
- Push changes to same branch
- Mark conversations as resolved

#### Don't Take It Personally
- Reviews are about code, not you
- Learn from feedback
- Ask questions if unsure

---

## Documentation Guidelines

### Code Comments

```go
// Good comments explain WHY, not WHAT
// Bad: Gets user by ID
// Good: Fetches user from cache first to reduce database load
func (s *Service) GetUser(id string) (*User, error) {
    // Check cache first
    if user := s.cache.Get(id); user != nil {
        return user, nil
    }
    
    // Cache miss - fetch from database
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache for future requests (1 hour TTL)
    s.cache.Set(id, user, time.Hour)
    
    return user, nil
}
```

### README Updates

Update service README when:
- Adding new API endpoints
- Changing configuration
- Adding dependencies
- Modifying deployment process

### API Documentation

Update API contracts when:
- Adding new endpoints
- Changing request/response formats
- Modifying error codes
- Updating authentication requirements

---

## Questions?

- Check existing documentation
- Ask in team channels
- Create an issue for clarification
- Reach out to team lead

---

Thank you for contributing! 🎉
