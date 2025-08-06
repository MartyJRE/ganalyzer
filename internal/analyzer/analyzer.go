package analyzer

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"ganalyzer/pkg/types"
)

type Analyzer struct {
	normalizer *NameNormalizer
	normalize  bool
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		normalizer: NewNameNormalizer(),
		normalize:  false,
	}
}

func NewAnalyzerWithNormalization(normalize bool) *Analyzer {
	return &Analyzer{
		normalizer: NewNameNormalizer(),
		normalize:  normalize,
	}
}

func (a *Analyzer) AnalyzeRepository(repoPath string) (*types.Repository, error) {
	repo := types.NewRepository(repoPath)

	if err := a.analyzeCommits(repo); err != nil {
		return nil, fmt.Errorf("failed to analyze commits in %s: %w", repoPath, err)
	}

	if err := a.analyzeLineChanges(repo); err != nil {
		return nil, fmt.Errorf("failed to analyze line changes in %s: %w", repoPath, err)
	}

	return repo, nil
}

func (a *Analyzer) getContributorKey(name string) string {
	if a.normalize {
		return a.normalizer.NormalizeName(name)
	}
	return name
}

func (a *Analyzer) analyzeCommits(repo *types.Repository) error {
	cmd := exec.Command("git", "-C", repo.Path, "shortlog", "-sn", "--all")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git shortlog failed: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	commitCountRegex := regexp.MustCompile(`^\s*(\d+)\s+(.+)$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := commitCountRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		count, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		authorName := strings.TrimSpace(matches[2])
		contributorKey := a.getContributorKey(authorName)

		if _, exists := repo.Contributors[contributorKey]; !exists {
			repo.Contributors[contributorKey] = &types.ContributorStats{
				Name:    authorName, // Keep original name for display
				Aliases: make([]string, 0),
			}
		}

		// Add this authorName as an alias if normalization is enabled and it's different from the stored name
		if a.normalize && authorName != repo.Contributors[contributorKey].Name {
			aliases := repo.Contributors[contributorKey].Aliases
			found := false
			for _, alias := range aliases {
				if alias == authorName {
					found = true
					break
				}
			}
			if !found {
				repo.Contributors[contributorKey].Aliases = append(repo.Contributors[contributorKey].Aliases, authorName)
			}
		}

		repo.Contributors[contributorKey].CommitCount += count
	}

	return scanner.Err()
}

func (a *Analyzer) analyzeLineChanges(repo *types.Repository) error {
	cmd := exec.Command("git", "-C", repo.Path, "log", "--all", "--format=%aN", "--numstat")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git log failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var currentAuthor string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.Contains(line, "\t") {
			currentAuthor = line
			continue
		}

		if currentAuthor == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}

		added, err1 := strconv.Atoi(parts[0])
		deleted, err2 := strconv.Atoi(parts[1])

		if err1 != nil || err2 != nil {
			continue
		}

		contributorKey := a.getContributorKey(currentAuthor)
		if _, exists := repo.Contributors[contributorKey]; !exists {
			repo.Contributors[contributorKey] = &types.ContributorStats{
				Name:    currentAuthor, // Keep original name for display
				Aliases: make([]string, 0),
			}
		}

		// Add this currentAuthor as an alias if normalization is enabled and it's different from the stored name
		if a.normalize && currentAuthor != repo.Contributors[contributorKey].Name {
			aliases := repo.Contributors[contributorKey].Aliases
			found := false
			for _, alias := range aliases {
				if alias == currentAuthor {
					found = true
					break
				}
			}
			if !found {
				repo.Contributors[contributorKey].Aliases = append(repo.Contributors[contributorKey].Aliases, currentAuthor)
			}
		}

		stats := repo.Contributors[contributorKey]
		stats.LinesAdded += added
		stats.LinesDeleted += deleted
		stats.LinesChanged += added + deleted
	}

	return nil
}
