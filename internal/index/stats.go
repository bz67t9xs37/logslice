package index

import (
	"fmt"
	"io"
	"time"
)

// Stats summarises an index built from a log file.
type Stats struct {
	// Entries is the total number of indexed lines.
	Entries int
	// First is the timestamp of the earliest indexed line.
	First time.Time
	// Last is the timestamp of the latest indexed line.
	Last time.Time
	// Span is the duration between First and Last.
	Span time.Duration
}

// Collect derives Stats from a non-empty slice of Entry values.
// It returns an error when entries is empty.
func Collect(entries []Entry) (Stats, error) {
	if len(entries) == 0 {
		return Stats{}, fmt.Errorf("index: cannot collect stats from empty entry slice")
	}
	s := Stats{
		Entries: len(entries),
		First:   entries[0].Timestamp,
		Last:    entries[len(entries)-1].Timestamp,
	}
	s.Span = s.Last.Sub(s.First)
	return s, nil
}

// Fprint writes a human-readable summary of s to w.
func (s Stats) Fprint(w io.Writer) {
	fmt.Fprintf(w, "entries : %d\n", s.Entries)
	fmt.Fprintf(w, "first   : %s\n", s.First.Format(time.RFC3339))
	fmt.Fprintf(w, "last    : %s\n", s.Last.Format(time.RFC3339))
	fmt.Fprintf(w, "span    : %s\n", s.Span)
}
