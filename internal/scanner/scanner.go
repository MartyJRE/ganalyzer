package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type Scanner struct {
	foundRepos []string
}

func NewScanner() *Scanner {
	return &Scanner{
		foundRepos: make([]string, 0),
	}
}

func (s *Scanner) ScanForRepositories(rootDir string) ([]string, error) {
	s.foundRepos = make([]string, 0)

	err := filepath.WalkDir(rootDir, s.walkFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory %s: %w", rootDir, err)
	}

	return s.filterNestedRepos(s.foundRepos), nil
}

func (s *Scanner) filterNestedRepos(repos []string) []string {
	if len(repos) <= 1 {
		return repos
	}

	filtered := make([]string, 0, len(repos))

	for i, repo := range repos {
		isNested := false
		for j, other := range repos {
			if i != j && strings.HasPrefix(repo, other+string(filepath.Separator)) {
				isNested = true
				break
			}
		}
		if !isNested {
			filtered = append(filtered, repo)
		}
	}

	return filtered
}

func (s *Scanner) walkFunc(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("Warning: cannot access %s: %v\n", path, err)
		return err
	}

	if !d.IsDir() {
		return nil
	}

	if d.Name() == ".git" {
		repoPath := filepath.Dir(path)
		s.foundRepos = append(s.foundRepos, repoPath)
		return fs.SkipDir
	}

	if shouldSkipDir(d.Name()) {
		return fs.SkipDir
	}

	return nil
}

func shouldSkipDir(dirName string) bool {
	skipDirs := []string{
		"node_modules",
		".vscode",
		".idea",
		"target",
		"build",
		"dist",
		".next",
		".nuxt",
		"vendor",
		"__pycache__",
		".cache",
		".DS_Store",
	}

	for _, skip := range skipDirs {
		if dirName == skip {
			return true
		}
	}

	return false
}
