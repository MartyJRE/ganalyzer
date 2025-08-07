package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"ganalyzer/internal/analyzer"
	"ganalyzer/internal/formatter"
	"ganalyzer/internal/scanner"
	"ganalyzer/pkg/types"
)

func main() {
	var config formatter.Config
	var showVersion bool

	flag.StringVar(&config.Directory, "dir", ".", "Directory to scan for Git repositories")
	flag.StringVar(&config.OutputFormat, "format", "table", "Output format: table, json, csv")
	flag.IntVar(&config.TopN, "top", 0, "Show only top N contributors (0 = all)")
	flag.StringVar(&config.SortBy, "sort", "commits", "Sort by: commits, lines, combined")
	flag.BoolVar(&config.NormalizeNames, "normalize", false, "Normalize contributor names (remove diacritics, punctuation, case differences)")
	flag.BoolVar(&config.ShowAliases, "aliases", false, "Show contributor aliases when normalization is enabled")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Println("ganalyzer version dev")
		os.Exit(0)
	}

	// Validate flag dependencies
	if config.ShowAliases && !config.NormalizeNames {
		fmt.Fprintf(os.Stderr, "Warning: -aliases flag requires -normalize flag to be effective\n")
	}

	absDir, err := filepath.Abs(config.Directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving directory path: %v\n", err)
		os.Exit(1)
	}
	config.Directory = absDir

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(config formatter.Config) error {
	repoScanner := scanner.NewScanner()
	repoAnalyzer := analyzer.NewAnalyzerWithNormalization(config.NormalizeNames)
	repoFormatter := formatter.NewFormatter()
	globalStats := types.NewGlobalStats()

	fmt.Fprintf(os.Stderr, "Scanning directory: %s\n", config.Directory)
	repos, err := repoScanner.ScanForRepositories(config.Directory)
	if err != nil {
		return fmt.Errorf("failed to scan for repositories: %w", err)
	}

	if len(repos) == 0 {
		fmt.Fprintf(os.Stderr, "No Git repositories found in %s\n", config.Directory)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d repositories, analyzing...\n", len(repos))

	for i, repoPath := range repos {
		fmt.Fprintf(os.Stderr, "Analyzing repository %d/%d: %s\n", i+1, len(repos), repoPath)

		repo, err := repoAnalyzer.AnalyzeRepository(repoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to analyze %s: %v\n", repoPath, err)
			continue
		}

		globalStats.AddRepository(repo)
	}

	return repoFormatter.Format(globalStats, config, os.Stdout)
}
