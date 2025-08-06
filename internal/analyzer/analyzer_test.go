package analyzer

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAnalyzer_Integration(t *testing.T) {
	if !hasGit() {
		t.Skip("Git not available, skipping integration test")
	}

	tempDir := createTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	analyzer := NewAnalyzer()
	repo, err := analyzer.AnalyzeRepository(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeRepository failed: %v", err)
	}

	if len(repo.Contributors) == 0 {
		t.Error("Expected at least one contributor")
	}

	if repo.Path != tempDir {
		t.Errorf("Expected repo path %s, got %s", tempDir, repo.Path)
	}

	testUser := repo.Contributors["Test User"]
	if testUser == nil {
		t.Error("Expected 'Test User' contributor not found")
	} else {
		if testUser.CommitCount == 0 {
			t.Error("Expected at least one commit for Test User")
		}
	}
}

func TestAnalyzer_NonexistentRepository(t *testing.T) {
	analyzer := NewAnalyzer()
	_, err := analyzer.AnalyzeRepository("/nonexistent/repo")
	if err == nil {
		t.Error("Expected error for nonexistent repository, got nil")
	}
}

func TestAnalyzer_NonGitDirectory(t *testing.T) {
	tempDir := t.TempDir()

	analyzer := NewAnalyzer()
	_, err := analyzer.AnalyzeRepository(tempDir)
	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
}

func hasGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func createTestGitRepo(t *testing.T) string {
	tempDir := t.TempDir()

	if err := runCmd(tempDir, "git", "init"); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	if err := runCmd(tempDir, "git", "config", "user.name", "Test User"); err != nil {
		t.Fatalf("git config user.name failed: %v", err)
	}

	if err := runCmd(tempDir, "git", "config", "user.email", "test@example.com"); err != nil {
		t.Fatalf("git config user.email failed: %v", err)
	}

	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Hello, World!\nSecond line\n"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	if err := runCmd(tempDir, "git", "add", "test.txt"); err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	if err := runCmd(tempDir, "git", "commit", "-m", "Initial commit"); err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	testFile2 := filepath.Join(tempDir, "test2.txt")
	if err := os.WriteFile(testFile2, []byte("Another file\nWith more content\n"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	if err := runCmd(tempDir, "git", "add", "test2.txt"); err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	if err := runCmd(tempDir, "git", "commit", "-m", "Second commit"); err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	return tempDir
}

func runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}
