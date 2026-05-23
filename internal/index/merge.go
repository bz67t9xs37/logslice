package index

import (
	"sort"
	"time"
)

// MergeEntries merges two sorted slices of Entry into a single sorted slice,
// deduplicating entries that share the same byte offset.
func MergeEntries(a, b []Entry) []Entry {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	combined := make([]Entry, 0, len(a)+len(b))
	combined = append(combined, a...)
	combined = append(combined, b...)

	sort.Slice(combined, func(i, j int) bool {
		if combined[i].Timestamp.Equal(combined[j].Timestamp) {
			return combined[i].Offset < combined[j].Offset
		}
		return combined[i].Timestamp.Before(combined[j].Timestamp)
	})

	return deduplicate(combined)
}

// deduplicate removes consecutive entries with identical offsets.
func deduplicate(entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}
	out := entries[:1]
	for _, e := range entries[1:] {
		if e.Offset != out[len(out)-1].Offset {
			out = append(out, e)
		}
	}
	return out
}

// Trim returns only the entries whose timestamps fall within [start, end] inclusive.
func Trim(entries []Entry, start, end time.Time) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !e.Timestamp.Before(start) && !e.Timestamp.After(end) {
			result = append(result, e)
		}
	}
	return result
}
