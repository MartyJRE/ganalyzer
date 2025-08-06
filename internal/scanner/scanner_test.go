package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_ScanForRepositories(t *testing.T) {
	tempDir := t.TempDir()

	repo1Path := filepath.Join(tempDir, "repo1")
	repo2Path := filepath.Join(tempDir, "subdir", "repo2")
	notRepoPath := filepath.Join(tempDir, "notrepo")
	nestedRepoPath := filepath.Join(repo1Path, "nested")

	if err := os.MkdirAll(filepath.Join(repo1Path, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create test repo1: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repo2Path, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create test repo2: %v", err)
	}
	if err := os.MkdirAll(notRepoPath, 0755); err != nil {
		t.Fatalf("Failed to create test notrepo: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(nestedRepoPath, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create nested repo: %v", err)
	}

	scanner := NewScanner()
	repos, err := scanner.ScanForRepositories(tempDir)
	if err != nil {
		t.Fatalf("ScanForRepositories failed: %v", err)
	}

	if len(repos) != 2 {
		t.Errorf("Expected 2 repositories, got %d: %v", len(repos), repos)
	}

	foundRepo1 := false
	foundRepo2 := false
	foundNested := false

	for _, repo := range repos {
		switch repo {
		case repo1Path:
			foundRepo1 = true
		case repo2Path:
			foundRepo2 = true
		case nestedRepoPath:
			foundNested = true
		}
	}

	if !foundRepo1 {
		t.Error("repo1 not found in results")
	}
	if !foundRepo2 {
		t.Error("repo2 not found in results")
	}
	if foundNested {
		t.Error("nested repo should not be found (first-level only)")
	}
}

func TestShouldSkipDir(t *testing.T) {
	tests := []struct {
		dirName  string
		expected bool
	}{
		{"node_modules", true},
		{".vscode", true},
		{".idea", true},
		{"target", true},
		{"build", true},
		{"dist", true},
		{".next", true},
		{".nuxt", true},
		{"vendor", true},
		{"__pycache__", true},
		{".cache", true},
		{".DS_Store", true},
		{"src", false},
		{"tests", false},
		{"docs", false},
		{"normal-dir", false},
	}

	for _, test := range tests {
		result := shouldSkipDir(test.dirName)
		if result != test.expected {
			t.Errorf("shouldSkipDir(%s) = %v, expected %v", test.dirName, result, test.expected)
		}
	}
}

func TestScanner_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	scanner := NewScanner()
	repos, err := scanner.ScanForRepositories(tempDir)
	if err != nil {
		t.Fatalf("ScanForRepositories failed: %v", err)
	}

	if len(repos) != 0 {
		t.Errorf("Expected 0 repositories in empty directory, got %d", len(repos))
	}
}

func TestScanner_NonexistentDirectory(t *testing.T) {
	scanner := NewScanner()
	_, err := scanner.ScanForRepositories("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
}
