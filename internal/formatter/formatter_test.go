package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"ganalyzer/pkg/types"
)

func TestFormatter_FormatTable(t *testing.T) {
	formatter := NewFormatter()
	stats := createTestGlobalStats()
	config := Config{
		OutputFormat: "table",
		SortBy:       "commits",
		TopN:         0,
	}

	var buf bytes.Buffer
	err := formatter.Format(stats, config, &buf)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "Git Repository Analysis") {
		t.Error("Expected table header not found")
	}

	if !strings.Contains(output, "Alice") {
		t.Error("Expected contributor 'Alice' not found in output")
	}

	if !strings.Contains(output, "Bob") {
		t.Error("Expected contributor 'Bob' not found in output")
	}

	if !strings.Contains(output, "test-repo") {
		t.Error("Expected repository name not found in output")
	}
}

func TestFormatter_FormatJSON(t *testing.T) {
	formatter := NewFormatter()
	stats := createTestGlobalStats()
	config := Config{
		OutputFormat: "json",
		SortBy:       "commits",
		TopN:         0,
	}

	var buf bytes.Buffer
	err := formatter.Format(stats, config, &buf)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	var result struct {
		Repositories []*types.Repository       `json:"repositories"`
		Contributors []*types.ContributorStats `json:"contributors"`
	}

	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if len(result.Contributors) != 2 {
		t.Errorf("Expected 2 contributors in JSON, got %d", len(result.Contributors))
	}

	if len(result.Repositories) != 1 {
		t.Errorf("Expected 1 repository in JSON, got %d", len(result.Repositories))
	}
}

func TestFormatter_FormatCSV(t *testing.T) {
	formatter := NewFormatter()
	stats := createTestGlobalStats()
	config := Config{
		OutputFormat: "csv",
		SortBy:       "commits",
		TopN:         0,
	}

	var buf bytes.Buffer
	err := formatter.Format(stats, config, &buf)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines (header + 2 contributors), got %d", len(lines))
	}

	header := lines[0]
	expectedHeaders := []string{"Name", "Email", "Commits", "Lines Added", "Lines Deleted", "Total Lines"}
	for _, expectedHeader := range expectedHeaders {
		if !strings.Contains(header, expectedHeader) {
			t.Errorf("Expected header '%s' not found in CSV header: %s", expectedHeader, header)
		}
	}

	if !strings.Contains(output, "Alice") || !strings.Contains(output, "Bob") {
		t.Error("Expected contributors not found in CSV output")
	}
}

func TestFormatter_UnsupportedFormat(t *testing.T) {
	formatter := NewFormatter()
	stats := createTestGlobalStats()
	config := Config{
		OutputFormat: "xml",
		SortBy:       "commits",
		TopN:         0,
	}

	var buf bytes.Buffer
	err := formatter.Format(stats, config, &buf)
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("Expected 'unsupported output format' error, got: %v", err)
	}
}

func TestFormatter_EmptyStats(t *testing.T) {
	formatter := NewFormatter()
	stats := types.NewGlobalStats()
	config := Config{
		OutputFormat: "table",
		SortBy:       "commits",
		TopN:         0,
	}

	var buf bytes.Buffer
	err := formatter.Format(stats, config, &buf)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No contributors found") {
		t.Error("Expected 'No contributors found' message for empty stats")
	}
}

func TestFormatter_CalculateNameWidthWithUnicode(t *testing.T) {
	formatter := NewFormatter()
	unicodeName := strings.Repeat("á", 21)
	contributors := []*types.ContributorStats{{Name: unicodeName}}
	width := formatter.calculateNameWidth(contributors, Config{})
	expected := utf8.RuneCountInString(unicodeName) + namePadding
	if width != expected {
		t.Fatalf("expected width %d, got %d", expected, width)
	}
}

func createTestGlobalStats() *types.GlobalStats {
	stats := types.NewGlobalStats()

	repo := types.NewRepository("/path/to/test-repo")
	repo.Contributors["Alice"] = &types.ContributorStats{
		Name:         "Alice",
		Email:        "alice@example.com",
		CommitCount:  10,
		LinesAdded:   100,
		LinesDeleted: 20,
		LinesChanged: 120,
	}
	repo.Contributors["Bob"] = &types.ContributorStats{
		Name:         "Bob",
		Email:        "bob@example.com",
		CommitCount:  5,
		LinesAdded:   50,
		LinesDeleted: 10,
		LinesChanged: 60,
	}

	stats.AddRepository(repo)
	return stats
}
