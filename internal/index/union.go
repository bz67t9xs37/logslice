// Package index provides log index building, caching, and querying utilities.
package index

import (
	"fmt"
	"os"
	"time"
)

// UnionResult holds the merged index entries and metadata about the source files.
type UnionResult struct {
	Entries []Entry
	Sources []string
}

// Union builds or loads indexes for each provided log file path and merges
// them into a single sorted, deduplicated entry list covering [start, end].
func Union(paths []string, start, end time.Time) (*UnionResult, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("index/union: no paths provided")
	}

	result := &UnionResult{
		Sources: make([]string, 0, len(paths)),
	}

	var merged []Entry

	for _, p := range paths {
		entries, err := loadOrBuild(p)
		if err != nil {
			return nil, fmt.Errorf("index/union: %s: %w", p, err)
		}
		trimmed := Trim(entries, start, end)
		if len(trimmed) == 0 {
			continue
		}
		merged = MergeEntries(merged, trimmed)
		result.Sources = append(result.Sources, p)
	}

	result.Entries = merged
	return result, nil
}

// loadOrBuild attempts to load a cached index for path; on miss it builds and
// saves a fresh one.
func loadOrBuild(path string) ([]Entry, error) {
	cachePath := CachePath(path)
	if entries, err := LoadCache(cachePath, path); err == nil {
		return entries, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entries, err := Build(f, nil)
	if err != nil {
		return nil, err
	}

	// Best-effort cache save — ignore errors.
	_ = SaveCache(cachePath, path, entries)
	return entries, nil
}
