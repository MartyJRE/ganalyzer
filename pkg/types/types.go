package types

import (
	"sort"
)

const (
	// Scoring weights for combined sorting
	commitsWeight = 10
	linesWeight   = 100
)

// Repository represents a Git repository with its contributor statistics
type Repository struct {
	Path         string
	Name         string
	Contributors map[string]*ContributorStats
}

// ContributorStats holds statistics for a single contributor
type ContributorStats struct {
	Name         string
	Email        string
	CommitCount  int
	LinesAdded   int
	LinesDeleted int
	LinesChanged int
	Aliases      []string
}

// GlobalStats aggregates contributor statistics across multiple repositories
type GlobalStats struct {
	Contributors map[string]*ContributorStats
	Repositories []*Repository
}

// NewGlobalStats creates a new GlobalStats instance
func NewGlobalStats() *GlobalStats {
	return &GlobalStats{
		Contributors: make(map[string]*ContributorStats),
		Repositories: make([]*Repository, 0),
	}
}

// AddRepository adds a repository's statistics to the global stats
func (gs *GlobalStats) AddRepository(repo *Repository) {
	gs.Repositories = append(gs.Repositories, repo)

	for name, stats := range repo.Contributors {
		if existing, exists := gs.Contributors[name]; exists {
			existing.CommitCount += stats.CommitCount
			existing.LinesAdded += stats.LinesAdded
			existing.LinesDeleted += stats.LinesDeleted
			existing.LinesChanged += stats.LinesChanged
			// Merge aliases, avoiding duplicates
			for _, alias := range stats.Aliases {
				found := false
				for _, existingAlias := range existing.Aliases {
					if existingAlias == alias {
						found = true
						break
					}
				}
				if !found {
					existing.Aliases = append(existing.Aliases, alias)
				}
			}
		} else {
			// Create a copy of aliases slice
			aliases := make([]string, len(stats.Aliases))
			copy(aliases, stats.Aliases)
			gs.Contributors[name] = &ContributorStats{
				Name:         stats.Name,
				Email:        stats.Email,
				CommitCount:  stats.CommitCount,
				LinesAdded:   stats.LinesAdded,
				LinesDeleted: stats.LinesDeleted,
				LinesChanged: stats.LinesChanged,
				Aliases:      aliases,
			}
		}
	}
}

// GetSortedContributors returns contributors sorted by the specified criteria
func (gs *GlobalStats) GetSortedContributors(sortBy string, topN int) []*ContributorStats {
	contributors := make([]*ContributorStats, 0, len(gs.Contributors))
	for _, stats := range gs.Contributors {
		contributors = append(contributors, stats)
	}

	switch sortBy {
	case "commits":
		sort.Slice(contributors, func(i, j int) bool {
			return contributors[i].CommitCount > contributors[j].CommitCount
		})
	case "lines":
		sort.Slice(contributors, func(i, j int) bool {
			return contributors[i].LinesChanged > contributors[j].LinesChanged
		})
	case "combined":
		sort.Slice(contributors, func(i, j int) bool {
			scoreI := contributors[i].CommitCount*commitsWeight + contributors[i].LinesChanged/linesWeight
			scoreJ := contributors[j].CommitCount*commitsWeight + contributors[j].LinesChanged/linesWeight
			return scoreI > scoreJ
		})
	default:
		sort.Slice(contributors, func(i, j int) bool {
			return contributors[i].CommitCount > contributors[j].CommitCount
		})
	}

	if topN > 0 && topN < len(contributors) {
		contributors = contributors[:topN]
	}

	return contributors
}

// NewRepository creates a new Repository instance for the given path
func NewRepository(path string) *Repository {
	return &Repository{
		Path:         path,
		Name:         extractRepoName(path),
		Contributors: make(map[string]*ContributorStats),
	}
}

func extractRepoName(path string) string {
	if path == "" {
		return "unknown"
	}

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}
