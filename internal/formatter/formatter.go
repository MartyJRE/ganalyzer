package formatter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"ganalyzer/pkg/types"
)

const (
	// Default minimum name column width
	minNameWidth = 20
	// Extra padding for name column
	namePadding = 2
)

// Config holds configuration options for formatting output
type Config struct {
	Directory      string
	OutputFormat   string
	TopN           int
	SortBy         string
	NormalizeNames bool
	ShowAliases    bool
}

// Formatter handles output formatting for analysis results
type Formatter struct{}

// NewFormatter creates a new Formatter instance
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Format outputs the analysis results in the specified format
func (f *Formatter) Format(stats *types.GlobalStats, config Config, writer io.Writer) error {
	contributors := stats.GetSortedContributors(config.SortBy, config.TopN)

	switch config.OutputFormat {
	case "json":
		return f.formatJSON(contributors, stats.Repositories, writer)
	case "csv":
		return f.formatCSV(contributors, config, writer)
	case "table":
		return f.formatTable(contributors, stats.Repositories, config, writer)
	default:
		return fmt.Errorf("unsupported output format: %s", config.OutputFormat)
	}
}

func (f *Formatter) formatTable(contributors []*types.ContributorStats, repos []*types.Repository, config Config, writer io.Writer) error {
	if err := f.writeHeader(writer, len(repos), repos); err != nil {
		return err
	}

	if len(contributors) == 0 {
		_, err := fmt.Fprintf(writer, "No contributors found.\n")
		return err
	}

	if err := f.writeContributorsHeader(writer); err != nil {
		return err
	}

	return f.writeContributors(writer, contributors, config)
}

func (f *Formatter) writeHeader(writer io.Writer, repoCount int, repos []*types.Repository) error {
	if _, err := fmt.Fprintf(writer, "Git Repository Analysis\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "======================\n\n"); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(writer, "Found %d repositories:\n", repoCount); err != nil {
		return err
	}
	for _, repo := range repos {
		if _, err := fmt.Fprintf(writer, "  - %s (%s)\n", repo.Name, repo.Path); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(writer, "\n")
	return err
}

func (f *Formatter) writeContributorsHeader(writer io.Writer) error {
	if _, err := fmt.Fprintf(writer, "Top Contributors:\n"); err != nil {
		return err
	}
	_, err := fmt.Fprintf(writer, "================\n\n")
	return err
}

func (f *Formatter) writeContributors(writer io.Writer, contributors []*types.ContributorStats, config Config) error {
	nameWidth := f.calculateNameWidth(contributors, config)
	format := fmt.Sprintf("%%-%ds %%8s %%10s %%10s %%12s\n", nameWidth)

	if err := f.writeTableHeader(writer, format, nameWidth); err != nil {
		return err
	}

	return f.writeContributorRows(writer, contributors, config, format)
}

func (f *Formatter) calculateNameWidth(contributors []*types.ContributorStats, config Config) int {
	nameWidth := minNameWidth
	for _, contributor := range contributors {
		name := f.formatContributorName(contributor, config)
		if len(name) > nameWidth {
			nameWidth = len(name)
		}
	}
	return nameWidth + namePadding
}

func (f *Formatter) formatContributorName(contributor *types.ContributorStats, config Config) string {
	if config.ShowAliases && config.NormalizeNames && len(contributor.Aliases) > 0 {
		return fmt.Sprintf("%s (aliases: %s)", contributor.Name, strings.Join(contributor.Aliases, ", "))
	}
	return contributor.Name
}

func (f *Formatter) writeTableHeader(writer io.Writer, format string, nameWidth int) error {
	if _, err := fmt.Fprintf(writer, format, "Name", "Commits", "Lines+", "Lines-", "Total Lines"); err != nil {
		return err
	}
	_, err := fmt.Fprintf(writer, format, strings.Repeat("-", nameWidth), "-------", "------", "------", "-----------")
	return err
}

func (f *Formatter) writeContributorRows(writer io.Writer, contributors []*types.ContributorStats, config Config, format string) error {
	for _, contributor := range contributors {
		name := f.formatContributorName(contributor, config)
		if _, err := fmt.Fprintf(writer, format,
			name,
			strconv.Itoa(contributor.CommitCount),
			strconv.Itoa(contributor.LinesAdded),
			strconv.Itoa(contributor.LinesDeleted),
			strconv.Itoa(contributor.LinesChanged),
		); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) formatJSON(contributors []*types.ContributorStats, repos []*types.Repository, writer io.Writer) error {
	data := struct {
		Repositories []*types.Repository       `json:"repositories"`
		Contributors []*types.ContributorStats `json:"contributors"`
	}{
		Repositories: repos,
		Contributors: contributors,
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *Formatter) formatCSV(contributors []*types.ContributorStats, config Config, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	headers := []string{"Name", "Email", "Commits", "Lines Added", "Lines Deleted", "Total Lines"}
	if config.ShowAliases && config.NormalizeNames {
		headers = append(headers, "Aliases")
	}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	for _, contributor := range contributors {
		record := []string{
			contributor.Name,
			contributor.Email,
			strconv.Itoa(contributor.CommitCount),
			strconv.Itoa(contributor.LinesAdded),
			strconv.Itoa(contributor.LinesDeleted),
			strconv.Itoa(contributor.LinesChanged),
		}
		if config.ShowAliases && config.NormalizeNames {
			record = append(record, strings.Join(contributor.Aliases, "; "))
		}
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	return nil
}
