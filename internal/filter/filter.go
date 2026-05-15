// Package filter provides log line filtering based on time range boundaries.
package filter

import (
	"time"

	"github.com/logslice/logslice/internal/timeparse"
)

// Result represents the outcome of evaluating a log line against a time range.
type Result int

const (
	// Before indicates the log line timestamp is before the target range.
	Before Result = iota
	// Inside indicates the log line timestamp falls within the target range.
	Inside
	// After indicates the log line timestamp is after the target range.
	After
	// Unparseable indicates no timestamp could be extracted from the line.
	Unparseable
)

// Filter evaluates log lines against a time range window.
type Filter struct {
	parser *timeparse.Parser
	start  time.Time
	end    time.Time
}

// New creates a Filter that accepts lines within [start, end].
func New(parser *timeparse.Parser, start, end time.Time) *Filter {
	return &Filter{
		parser: parser,
		start:  start,
		end:    end,
	}
}

// Evaluate checks a raw log line and returns its Result relative to the
// configured time range. Lines that cannot be parsed are returned as
// Unparseable so callers can decide how to handle them (e.g. carry-forward).
func (f *Filter) Evaluate(line string) Result {
	t, err := f.parser.Parse(line)
	if err != nil {
		return Unparseable
	}

	switch {
	case t.Before(f.start):
		return Before
	case t.After(f.end):
		return After
	default:
		return Inside
	}
}

// InRange returns true when the line's timestamp falls within [start, end].
func (f *Filter) InRange(line string) bool {
	return f.Evaluate(line) == Inside
}
