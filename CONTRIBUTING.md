# Contributing to ganalyzer

Thank you for your interest in contributing to ganalyzer! We welcome contributions from everyone, whether you're fixing a bug, adding a feature, improving documentation, or just asking questions.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Issue Guidelines](#issue-guidelines)
- [Feature Requests](#feature-requests)
- [Questions and Support](#questions-and-support)

## 🤝 Code of Conduct

This project adheres to a code of conduct adapted from the [Contributor Covenant](https://www.contributor-covenant.org/). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- Git
- Make (optional, but recommended)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/ganalyzer.git
   cd ganalyzer
   ```

3. Add the original repository as upstream:
   ```bash
   git remote add upstream https://github.com/MartyJRE/ganalyzer.git
   ```

## 🛠 Development Setup

### Install Dependencies

ganalyzer uses only Go standard library, so no external dependencies are needed:

```bash
# Verify Go installation
go version

# Build the project
make build

# Or build manually
go build -o build/ganalyzer ./cmd/ganalyzer
```

### Development Workflow

```bash
# Full development workflow (recommended)
make dev

# Individual steps
make fmt      # Format code
make vet      # Run go vet
make test     # Run tests
make build    # Build binary
```

### Verify Setup

```bash
# Run the binary
./build/ganalyzer -dir . -top 5

# Should output contributor analysis of the ganalyzer repository
```

## ✏️ Making Changes

### Create a Branch

Always create a new branch for your changes:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### Branch Naming Convention

- `feature/description` - for new features
- `fix/description` - for bug fixes
- `docs/description` - for documentation changes
- `refactor/description` - for code refactoring
- `test/description` - for test improvements

### Making Commits

Write clear, descriptive commit messages:

```bash
# Good commit messages
git commit -m "Add alias tracking for contributor normalization"
git commit -m "Fix commit count aggregation bug in analyzer"
git commit -m "Update README with normalization examples"

# Less helpful
git commit -m "Fix bug"
git commit -m "Update code"
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/analyzer/

# Run with verbose output
go test -v ./...
```

### Writing Tests

- Place tests in the same package as the code being tested
- Use descriptive test names: `TestAnalyzer_NormalizesContributorNames`
- Test both success and error cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies (Git commands) when appropriate

Example test structure:
```go
func TestAnalyzer_AnalyzeRepository(t *testing.T) {
    tests := []struct {
        name           string
        repoPath       string
        expectedError  bool
        expectedCount  int
    }{
        {
            name:          "valid repository",
            repoPath:      "/path/to/test/repo",
            expectedError: false,
            expectedCount: 5,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Coverage

Aim for reasonable test coverage:
- New features should include tests
- Bug fixes should include regression tests
- Critical paths should have comprehensive coverage

Check coverage with:
```bash
make coverage
```

## 📤 Submitting Changes

### Before Submitting

1. **Sync with upstream**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run the full test suite**:
   ```bash
   make dev
   ```

3. **Test manually** with various scenarios:
   ```bash
   # Test different flags and edge cases
   ./build/ganalyzer -normalize -aliases
   ./build/ganalyzer -format json
   ./build/ganalyzer -dir /nonexistent/path
   ```

### Pull Request Process

1. **Push your branch**:
   ```bash
   git push origin your-branch-name
   ```

2. **Create a Pull Request** on GitHub with:
   - Clear title describing the change
   - Detailed description of what changed and why
   - Link to any related issues
   - Screenshots/examples if applicable

3. **PR Template** (use this structure):
   ```markdown
   ## Summary
   Brief description of changes
   
   ## Changes
   - List of specific changes made
   - Another change
   
   ## Testing
   - [ ] Unit tests pass
   - [ ] Manual testing completed
   - [ ] Edge cases considered
   
   ## Related Issues
   Fixes #123
   ```

4. **Address Review Feedback** promptly and respectfully

### PR Review Criteria

Your PR will be evaluated on:
- **Functionality**: Does it work as intended?
- **Code Quality**: Is it readable, maintainable, and well-structured?
- **Testing**: Are there adequate tests?
- **Documentation**: Are changes documented appropriately?
- **Performance**: Does it maintain or improve performance?
- **Compatibility**: Does it break existing functionality?

## 📝 Code Style

### Go Style Guidelines

We follow standard Go conventions:

- **Formatting**: Use `gofmt` (run `make fmt`)
- **Linting**: Use `go vet` (run `make vet`)
- **Naming**: Follow Go naming conventions
  - Use `camelCase` for variables and functions
  - Use `PascalCase` for exported functions and types
  - Use descriptive names
- **Comments**: 
  - Public functions must have comments
  - Comments should explain "why", not "what"
- **Error handling**: Always handle errors appropriately

### Project-Specific Guidelines

- **Package structure**: Keep packages focused and cohesive
- **Interfaces**: Keep them small and focused
- **Configuration**: Use struct fields rather than global variables
- **Constants**: Group related constants together
- **File organization**: One main concept per file

### Example Code Style

```go
// Good
func (a *Analyzer) AnalyzeRepository(repoPath string) (*types.Repository, error) {
    if repoPath == "" {
        return nil, fmt.Errorf("repository path cannot be empty")
    }
    
    repo := types.NewRepository(repoPath)
    
    if err := a.analyzeCommits(repo); err != nil {
        return nil, fmt.Errorf("failed to analyze commits: %w", err)
    }
    
    return repo, nil
}

// Less preferred
func (a *Analyzer) analyze(path string) (*types.Repository, error) {
    r := &types.Repository{Path: path}
    a.analyzeCommits(r) // Error not handled
    return r, nil
}
```

## 🐛 Issue Guidelines

### Reporting Bugs

Use the bug report template and include:

1. **Environment**:
   - OS and version
   - Go version
   - ganalyzer version

2. **Expected vs Actual Behavior**:
   - What you expected to happen
   - What actually happened

3. **Reproduction Steps**:
   - Minimal steps to reproduce
   - Command line used
   - Sample repository structure if relevant

4. **Additional Context**:
   - Error messages
   - Logs or output
   - Screenshots if applicable

### Example Bug Report

```markdown
## Bug Description
Normalization incorrectly groups different people with similar names

## Environment
- OS: macOS 14.0
- Go: 1.24.0
- ganalyzer: v1.0.0

## Steps to Reproduce
1. Run `./ganalyzer -normalize -dir /path/to/repo`
2. Observe output shows "John Smith" and "Jane Smith" as same contributor

## Expected Behavior
Different people should remain separate contributors

## Actual Behavior
Names are incorrectly merged due to surname matching

## Additional Context
Repository has two developers: John Smith <john@email.com> and Jane Smith <jane@email.com>
```

## 💡 Feature Requests

### Proposing Features

1. **Check existing issues** first to avoid duplicates
2. **Describe the problem** the feature would solve
3. **Propose a solution** with specific details
4. **Consider alternatives** and trade-offs
5. **Assess impact** on existing functionality

### Feature Request Template

```markdown
## Problem Statement
What problem does this feature solve?

## Proposed Solution
Detailed description of the proposed feature

## Alternatives Considered
Other approaches considered and why they were not chosen

## Additional Context
- Use cases
- Examples
- References
```

## ❓ Questions and Support

### Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions, ideas, and general discussion
- **Code Review**: Ask questions in PR comments

### Documentation

- **README.md**: Usage instructions and examples
- **CLAUDE.md**: Development guidelines and project structure
- **Code comments**: Implementation details

## 🏆 Recognition

Contributors will be recognized in:
- Git commit history
- GitHub contributors list
- Future changelog and release notes

Thank you for contributing to ganalyzer! Your efforts help make this tool better for everyone. 🎉