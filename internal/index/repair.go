package index

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// RepairResult holds the outcome of a repair operation.
type RepairResult struct {
	OriginalCount int
	RepairedCount int
	RemovedZero   int
	RemovedDups   int
	Reordered     bool
}

// Fprint writes a human-readable summary of the repair result to w.
func (r RepairResult) Fprint(w io.Writer) {
	fmt.Fprintf(w, "repair: original=%d repaired=%d removed_zero=%d removed_dups=%d reordered=%v\n",
		r.OriginalCount, r.RepairedCount, r.RemovedZero, r.RemovedDups, r.Reordered)
}

// Repair attempts to fix a slice of entries by removing zero timestamps,
// deduplicating by offset, and sorting by timestamp.
// It returns the cleaned entries and a RepairResult describing what changed.
func Repair(entries []Entry) ([]Entry, RepairResult) {
	res := RepairResult{OriginalCount: len(entries)}

	// Remove zero timestamps.
	filtered := entries[:0:0]
	for _, e := range entries {
		if e.Timestamp.IsZero() {
			res.RemovedZero++
			continue
		}
		filtered = append(filtered, e)
	}

	// Deduplicate by offset, keeping earliest timestamp.
	seen := make(map[int64]time.Time, len(filtered))
	for _, e := range filtered {
		if t, ok := seen[e.Offset]; !ok || e.Timestamp.Before(t) {
			seen[e.Offset] = e.Timestamp
		}
	}
	deduped := make([]Entry, 0, len(seen))
	for offset, ts := range seen {
		deduped = append(deduped, Entry{Timestamp: ts, Offset: offset})
	}
	res.RemovedDups = len(filtered) - len(deduped)

	// Detect if sorting is needed.
	isSorted := sort.SliceIsSorted(deduped, func(i, j int) bool {
		return deduped[i].Timestamp.Before(deduped[j].Timestamp)
	})
	if !isSorted {
		sort.Slice(deduped, func(i, j int) bool {
			return deduped[i].Timestamp.Before(deduped[j].Timestamp)
		})
		res.Reordered = true
	}

	res.RepairedCount = len(deduped)
	return deduped, res
}
