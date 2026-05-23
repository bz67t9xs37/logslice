package index

import (
	"sort"
	"time"
)

// Entry represents a single indexed log line with its timestamp and byte offset.
type Entry struct {
	Timestamp time.Time
	Offset    int64
}

// Entries is a sorted slice of Entry values.
type Entries []Entry

// FindRange returns the byte offsets [start, end) that cover all log lines
// whose timestamps fall within [from, to] (inclusive). If no entries fall
// within the range, ok is false.
func (e Entries) FindRange(from, to time.Time) (startOffset, endOffset int64, ok bool) {
	if len(e) == 0 {
		return 0, 0, false
	}

	// Find first entry >= from
	lo := sort.Search(len(e), func(i int) bool {
		return !e[i].Timestamp.Before(from)
	})

	if lo >= len(e) {
		return 0, 0, false
	}

	// Find first entry > to
	hi := sort.Search(len(e), func(i int) bool {
		return e[i].Timestamp.After(to)
	})

	if hi == lo {
		return 0, 0, false
	}

	startOffset = e[lo].Offset

	// endOffset: if hi is within bounds use that entry's offset,
	// otherwise signal EOF with -1.
	if hi < len(e) {
		endOffset = e[hi].Offset
	} else {
		endOffset = -1 // caller should read until EOF
	}

	return startOffset, endOffset, true
}

// Len implements sort.Interface.
func (e Entries) Len() int { return len(e) }

// Less implements sort.Interface.
func (e Entries) Less(i, j int) bool { return e[i].Timestamp.Before(e[j].Timestamp) }

// Swap implements sort.Interface.
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
