package types

import (
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository("/path/to/my-repo")

	if repo.Path != "/path/to/my-repo" {
		t.Errorf("Expected path '/path/to/my-repo', got '%s'", repo.Path)
	}

	if repo.Name != "my-repo" {
		t.Errorf("Expected name 'my-repo', got '%s'", repo.Name)
	}

	if repo.Contributors == nil {
		t.Error("Contributors map should be initialized")
	}
}

func TestExtractRepoName(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/my-repo", "my-repo"},
		{"C:\\Users\\test\\project", "project"},
		{"simple-name", "simple-name"},
		{"", "unknown"},
		{"/", ""},
		{"path/", ""},
	}

	for _, test := range tests {
		result := extractRepoName(test.path)
		if result != test.expected {
			t.Errorf("extractRepoName(%s) = %s, expected %s", test.path, result, test.expected)
		}
	}
}

func TestGlobalStats_AddRepository(t *testing.T) {
	gs := NewGlobalStats()

	repo1 := NewRepository("/repo1")
	repo1.Contributors["alice"] = &ContributorStats{
		Name:         "alice",
		Email:        "alice@example.com",
		CommitCount:  10,
		LinesAdded:   100,
		LinesDeleted: 20,
		LinesChanged: 120,
	}
	repo1.Contributors["bob"] = &ContributorStats{
		Name:         "bob",
		Email:        "bob@example.com",
		CommitCount:  5,
		LinesAdded:   50,
		LinesDeleted: 10,
		LinesChanged: 60,
	}

	gs.AddRepository(repo1)

	if len(gs.Repositories) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(gs.Repositories))
	}

	if len(gs.Contributors) != 2 {
		t.Errorf("Expected 2 contributors, got %d", len(gs.Contributors))
	}

	alice := gs.Contributors["alice"]
	if alice.CommitCount != 10 || alice.LinesChanged != 120 {
		t.Errorf("Alice stats incorrect: commits=%d, lines=%d", alice.CommitCount, alice.LinesChanged)
	}

	repo2 := NewRepository("/repo2")
	repo2.Contributors["alice"] = &ContributorStats{
		Name:         "alice",
		Email:        "alice@example.com",
		CommitCount:  3,
		LinesAdded:   30,
		LinesDeleted: 5,
		LinesChanged: 35,
	}

	gs.AddRepository(repo2)

	if len(gs.Contributors) != 2 {
		t.Errorf("Expected 2 contributors after second repo, got %d", len(gs.Contributors))
	}

	alice = gs.Contributors["alice"]
	if alice.CommitCount != 13 || alice.LinesChanged != 155 {
		t.Errorf("Alice aggregated stats incorrect: commits=%d, lines=%d", alice.CommitCount, alice.LinesChanged)
	}
}

func TestGlobalStats_GetSortedContributors(t *testing.T) {
	gs := NewGlobalStats()
	gs.Contributors["alice"] = &ContributorStats{
		Name:         "alice",
		CommitCount:  10,
		LinesChanged: 100,
	}
	gs.Contributors["bob"] = &ContributorStats{
		Name:         "bob",
		CommitCount:  5,
		LinesChanged: 200,
	}
	gs.Contributors["charlie"] = &ContributorStats{
		Name:         "charlie",
		CommitCount:  15,
		LinesChanged: 50,
	}

	t.Run("sort by commits", func(t *testing.T) {
		sorted := gs.GetSortedContributors("commits", 0)
		if len(sorted) != 3 {
			t.Errorf("Expected 3 contributors, got %d", len(sorted))
		}
		if sorted[0].Name != "charlie" || sorted[0].CommitCount != 15 {
			t.Errorf("Expected charlie first with 15 commits, got %s with %d", sorted[0].Name, sorted[0].CommitCount)
		}
	})

	t.Run("sort by lines", func(t *testing.T) {
		sorted := gs.GetSortedContributors("lines", 0)
		if sorted[0].Name != "bob" || sorted[0].LinesChanged != 200 {
			t.Errorf("Expected bob first with 200 lines, got %s with %d", sorted[0].Name, sorted[0].LinesChanged)
		}
	})

	t.Run("limit results", func(t *testing.T) {
		sorted := gs.GetSortedContributors("commits", 2)
		if len(sorted) != 2 {
			t.Errorf("Expected 2 contributors with limit, got %d", len(sorted))
		}
	})
}
