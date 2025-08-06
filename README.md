# ganalyzer

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-007d9c.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Analyze Git repositories at scale.** Recursively scan directories, rank contributors by commits and code changes, normalize author names with alias tracking, and export results in multiple formats. Perfect for understanding team contributions across large codebases.

## 🔍 Features

- **Recursive repository discovery** - Automatically finds all Git repositories in directory trees
- **Multi-metric contributor ranking** - Sort by commits, lines changed, or combined score
- **Smart name normalization** - Handles variations like "John Doe", "john.doe", and "J. Doe" as one contributor
- **Alias tracking** - See all name variants used by each contributor
- **Multiple export formats** - Table (default), JSON, and CSV output
- **High performance** - Built with Go, uses only standard library
- **Comprehensive filtering** - Skips common build/cache directories automatically
- **Progress reporting** - Real-time feedback during analysis of large directory trees

## 📊 Use Cases

- **Team contribution analysis** across microservices and multi-repo projects
- **Code ownership assessment** for large, distributed codebases
- **Migration impact analysis** when moving between version control systems
- **Performance review data** collection with detailed metrics
- **Open source project insights** for maintainers and contributors

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/MartyJRE/ganalyzer.git
cd ganalyzer

# Build the binary
make build

# Or build directly with Go
go build -o build/ganalyzer ./cmd/ganalyzer
```

### Basic Usage

```bash
# Analyze current directory
./build/ganalyzer

# Analyze specific directory with normalized names
./build/ganalyzer -dir /path/to/projects -normalize

# Show top 10 contributors with aliases
./build/ganalyzer -normalize -aliases -top 10

# Export to CSV
./build/ganalyzer -format csv > contributors.csv

# Sort by lines changed instead of commits
./build/ganalyzer -sort lines
```

## 📖 Usage Examples

### Basic Analysis
```bash
./ganalyzer -dir ~/projects
```

### Advanced Analysis with Normalization
```bash
./ganalyzer -dir ~/projects -normalize -aliases -top 20 -sort combined
```

### Export Options
```bash
# JSON export
./ganalyzer -format json > analysis.json

# CSV export with aliases
./ganalyzer -normalize -aliases -format csv > team_analysis.csv
```

## 🛠 Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-dir` | Directory to scan for Git repositories | `.` (current) |
| `-format` | Output format: `table`, `json`, `csv` | `table` |
| `-top` | Show only top N contributors (0 = all) | `0` |
| `-sort` | Sort by: `commits`, `lines`, `combined` | `commits` |
| `-normalize` | Normalize contributor names | `false` |
| `-aliases` | Show contributor aliases (requires `-normalize`) | `false` |

## 💡 Name Normalization

The normalization feature intelligently groups contributors who appear under different names:

**Before normalization:**
```
Martin Pražák          1,203 commits
Martin Prazak            944 commits  
martin.prazak             55 commits
```

**After normalization with aliases:**
```
Martin Pražák (aliases: Martin Prazak, martin.prazak)    2,202 commits
```

This handles common variations:
- **Diacritics**: "José García" ↔ "Jose Garcia"
- **Punctuation**: "john.doe" ↔ "john doe" ↔ "johndoe"
- **Case differences**: "John Smith" ↔ "john smith"
- **Name order**: "Smith, John" ↔ "John Smith"

## 📁 Output Formats

### Table Format (Default)
```
Top Contributors:
================

Name                                 Commits     Lines+     Lines-  Total Lines
-----------------------------------  -------     ------     ------  -----------
Martin Pražák (aliases: Martin P...)    4402    7564560   11958531     19523091
Michal Mozik                             3774   14980545   11122811     26103356
```

### JSON Format
```json
{
  "repositories": [
    {
      "path": "/path/to/repo",
      "name": "my-project"
    }
  ],
  "contributors": [
    {
      "name": "Martin Pražák",
      "commits": 4402,
      "lines_added": 7564560,
      "lines_deleted": 11958531,
      "lines_changed": 19523091,
      "aliases": ["Martin Prazak", "martin.prazak"]
    }
  ]
}
```

### CSV Format
```csv
Name,Email,Commits,Lines Added,Lines Deleted,Total Lines,Aliases
Martin Pražák,,4402,7564560,11958531,19523091,"Martin Prazak; martin.prazak"
```

## 🏗 Development

### Prerequisites
- Go 1.24 or higher
- Git

### Building
```bash
# Development workflow (format, vet, test, build)
make dev

# Individual commands
make build    # Build executable
make test     # Run tests
make coverage # Run tests with coverage
make clean    # Clean build artifacts
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

### Project Structure
```
ganalyzer/
├── cmd/ganalyzer/          # Application entry point
├── internal/               # Private application code
│   ├── analyzer/           # Git analysis logic
│   ├── scanner/            # Repository discovery
│   └── formatter/          # Output formatting
├── pkg/types/              # Shared data types
├── build/                  # Build artifacts
└── Makefile               # Build automation
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on:
- How to submit issues and feature requests
- Development setup and workflow
- Code style guidelines
- Testing requirements
- Pull request process

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔧 Architecture

ganalyzer follows clean architecture principles:

- **Scanner** - Discovers Git repositories using `filepath.WalkDir`
- **Analyzer** - Executes Git commands and extracts contributor data
- **Formatter** - Handles multiple output formats (table, JSON, CSV)
- **Types** - Shared data structures for repositories and contributor statistics

The tool only scans first-level repositories (doesn't recurse into found Git repositories) to avoid double-counting contributions and improve performance.

## 🚨 Known Limitations

- Requires Git to be installed and accessible in PATH
- Analyzes only Git repositories (no SVN, Mercurial, etc.)
- Name normalization works best with Latin scripts
- Large repositories may take significant time to analyze
- Binary files are handled but may affect line count accuracy

## 💬 Support

- **Issues**: [GitHub Issues](https://github.com/MartyJRE/ganalyzer/issues)
- **Discussions**: [GitHub Discussions](https://github.com/MartyJRE/ganalyzer/discussions)
- **Security**: Please report security issues privately via email

---

**Built with ❤️ using Go and targeting real-world enterprise development workflows.**