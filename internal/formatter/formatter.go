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

type Config struct {
	Directory       string
	OutputFormat    string
	TopN            int
	SortBy          string
	NormalizeNames  bool
	ShowAliases     bool
}

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

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
	fmt.Fprintf(writer, "Git Repository Analysis\n")
	fmt.Fprintf(writer, "======================\n\n")

	fmt.Fprintf(writer, "Found %d repositories:\n", len(repos))
	for _, repo := range repos {
		fmt.Fprintf(writer, "  - %s (%s)\n", repo.Name, repo.Path)
	}
	fmt.Fprintf(writer, "\n")

	if len(contributors) == 0 {
		fmt.Fprintf(writer, "No contributors found.\n")
		return nil
	}

	fmt.Fprintf(writer, "Top Contributors:\n")
	fmt.Fprintf(writer, "================\n\n")

	nameWidth := 20
	for _, contributor := range contributors {
		name := contributor.Name
		if config.ShowAliases && config.NormalizeNames && len(contributor.Aliases) > 0 {
			name = fmt.Sprintf("%s (aliases: %s)", contributor.Name, strings.Join(contributor.Aliases, ", "))
		}
		if len(name) > nameWidth {
			nameWidth = len(name)
		}
	}
	nameWidth += 2

	format := fmt.Sprintf("%%-%ds %%8s %%10s %%10s %%12s\n", nameWidth)
	fmt.Fprintf(writer, format, "Name", "Commits", "Lines+", "Lines-", "Total Lines")
	fmt.Fprintf(writer, format, strings.Repeat("-", nameWidth), "-------", "------", "------", "-----------")

	for _, contributor := range contributors {
		name := contributor.Name
		if config.ShowAliases && config.NormalizeNames && len(contributor.Aliases) > 0 {
			name = fmt.Sprintf("%s (aliases: %s)", contributor.Name, strings.Join(contributor.Aliases, ", "))
		}
		
		fmt.Fprintf(writer, format,
			name,
			strconv.Itoa(contributor.CommitCount),
			strconv.Itoa(contributor.LinesAdded),
			strconv.Itoa(contributor.LinesDeleted),
			strconv.Itoa(contributor.LinesChanged),
		)
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
