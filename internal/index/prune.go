package index

import (
	"time"
)

// PruneOptions controls how entries are removed from an index.
type PruneOptions struct {
	// OlderThan removes entries whose timestamp is before Now-OlderThan.
	OlderThan time.Duration
	// MaxEntries keeps only the most recent N entries. Zero means no limit.
	MaxEntries int
}

// Prune removes entries from entries according to opts and returns the
// filtered slice. The input slice is not modified.
func Prune(entries []Entry, opts PruneOptions) []Entry {
	if len(entries) == 0 {
		return entries
	}

	result := make([]Entry, 0, len(entries))

	var cutoff time.Time
	if opts.OlderThan > 0 {
		cutoff = time.Now().Add(-opts.OlderThan)
	}

	for _, e := range entries {
		if !cutoff.IsZero() && e.Timestamp.Before(cutoff) {
			continue
		}
		result = append(result, e)
	}

	if opts.MaxEntries > 0 && len(result) > opts.MaxEntries {
		// Keep the most recent entries (entries are assumed sorted ascending).
		result = result[len(result)-opts.MaxEntries:]
	}

	return result
}
