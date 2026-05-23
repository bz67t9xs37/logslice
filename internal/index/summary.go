package index

import (
	"fmt"
	"io"
	"time"
)

// Summary holds human-readable statistics derived from an index.
type Summary struct {
	FileCount    int
	EntryCount   int
	EarliestTime time.Time
	LatestTime   time.Time
	Span         time.Duration
	Density      float64 // entries per minute
}

// Summarize builds a Summary from a slice of Entry values across one or more
// logical files. fileCount is the number of source log files that contributed
// to entries.
func Summarize(entries []Entry, fileCount int) Summary {
	if len(entries) == 0 {
		return Summary{FileCount: fileCount}
	}

	earliestIdx, latestIdx := 0, 0
	for i, e := range entries {
		if e.Timestamp.Before(entries[earliestIdx].Timestamp) {
			earliestIdx = i
		}
		if e.Timestamp.After(entries[latestIdx].Timestamp) {
			latestIdx = i
		}
	}

	earliestTime := entries[earliestIdx].Timestamp
	latestTime := entries[latestIdx].Timestamp
	span := latestTime.Sub(earliestTime)

	var density float64
	if minutes := span.Minutes(); minutes > 0 {
		density = float64(len(entries)) / minutes
	}

	return Summary{
		FileCount:    fileCount,
		EntryCount:   len(entries),
		EarliestTime: earliestTime,
		LatestTime:   latestTime,
		Span:         span,
		Density:      density,
	}
}

// Fprint writes a human-readable summary to w.
func (s Summary) Fprint(w io.Writer) {
	fmt.Fprintf(w, "Files        : %d\n", s.FileCount)
	fmt.Fprintf(w, "Entries      : %d\n", s.EntryCount)
	if s.EntryCount == 0 {
		fmt.Fprintln(w, "No entries indexed.")
		return
	}
	fmt.Fprintf(w, "Earliest     : %s\n", s.EarliestTime.Format(time.RFC3339))
	fmt.Fprintf(w, "Latest       : %s\n", s.LatestTime.Format(time.RFC3339))
	fmt.Fprintf(w, "Span         : %s\n", s.Span.Round(time.Second))
	fmt.Fprintf(w, "Density      : %.2f entries/min\n", s.Density)
}
