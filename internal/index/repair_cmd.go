package index

import (
	"fmt"
	"os"
)

// RepairOptions controls behaviour of RunRepair.
type RepairOptions struct {
	// Verbose prints the RepairResult to stderr when true.
	Verbose bool
}

// DefaultRepairOptions returns sensible defaults.
func DefaultRepairOptions() RepairOptions {
	return RepairOptions{Verbose: false}
}

// RepairFile loads the index cache for logPath, repairs it, and saves it back.
// It returns the RepairResult so callers can inspect what changed.
func RepairFile(logPath string, opts RepairOptions) (RepairResult, error) {
	cachePath := CachePath(logPath)

	entries, err := LoadCache(cachePath, logPath)
	if err != nil {
		return RepairResult{}, fmt.Errorf("repair: load cache: %w", err)
	}

	repaired, result := Repair(entries)

	if err := SaveCache(cachePath, repaired); err != nil {
		return result, fmt.Errorf("repair: save cache: %w", err)
	}

	if opts.Verbose {
		result.Fprint(os.Stderr)
	}
	return result, nil
}

// RunRepair is a convenience wrapper that prints results to stdout.
func RunRepair(logPath string) error {
	result, err := RepairFile(logPath, RepairOptions{Verbose: false})
	if err != nil {
		return err
	}
	result.Fprint(os.Stdout)
	return nil
}
