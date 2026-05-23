package index

import (
	"fmt"
	"os"
	"time"
)

// CompactFile loads (or builds) the index for logPath, compacts it with opts,
// saves the result back to the cache, and returns the compacted entries along
// with a summary of what changed.
//
// This is a convenience wrapper intended for CLI use.
func CompactFile(logPath string, opts CompactOptions) ([]Entry, CompactResult, error) {
	entries, err := loadOrBuild(logPath)
	if err != nil {
		return nil, CompactResult{}, fmt.Errorf("compact: load index for %q: %w", logPath, err)
	}

	compacted, result := Compact(entries, opts)

	cachePath := CachePath(logPath)
	if err := SaveCache(cachePath, compacted); err != nil {
		return nil, result, fmt.Errorf("compact: save cache %q: %w", cachePath, err)
	}

	return compacted, result, nil
}

// DefaultCompactOptions returns sensible defaults for log compaction:
// 1-minute buckets with no hard entry cap.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{
		BucketSize: time.Minute,
		MaxEntries: 0,
	}
}

// RunCompact is the entry-point called by the CLI sub-command.
// It prints a human-readable result to stdout.
func RunCompact(logPath string, opts CompactOptions) error {
	_, result, err := CompactFile(logPath, opts)
	if err != nil {
		return err
	}
	result.Fprint(os.Stdout)
	return nil
}
