package filter_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/logslice/internal/timeparse"
)

func newParser(t *testing.T) *timeparse.Parser {
	t.Helper()
	p, err := timeparse.New()
	if err != nil {
		t.Fatalf("timeparse.New: %v", err)
	}
	return p
}

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("mustTime: %v", err)
	}
	return parsed
}

func TestEvaluate_Inside(t *testing.T) {
	p := newParser(t)
	start := mustTime(t, "2024-01-01T10:00:00Z")
	end := mustTime(t, "2024-01-01T11:00:00Z")
	f := filter.New(p, start, end)

	line := "2024-01-01T10:30:00Z INFO something happened"
	if got := f.Evaluate(line); got != filter.Inside {
		t.Errorf("expected Inside, got %v", got)
	}
}

func TestEvaluate_Before(t *testing.T) {
	p := newParser(t)
	start := mustTime(t, "2024-01-01T10:00:00Z")
	end := mustTime(t, "2024-01-01T11:00:00Z")
	f := filter.New(p, start, end)

	line := "2024-01-01T09:00:00Z INFO early line"
	if got := f.Evaluate(line); got != filter.Before {
		t.Errorf("expected Before, got %v", got)
	}
}

func TestEvaluate_After(t *testing.T) {
	p := newParser(t)
	start := mustTime(t, "2024-01-01T10:00:00Z")
	end := mustTime(t, "2024-01-01T11:00:00Z")
	f := filter.New(p, start, end)

	line := "2024-01-01T12:00:00Z INFO late line"
	if got := f.Evaluate(line); got != filter.After {
		t.Errorf("expected After, got %v", got)
	}
}

func TestEvaluate_Unparseable(t *testing.T) {
	p := newParser(t)
	start := mustTime(t, "2024-01-01T10:00:00Z")
	end := mustTime(t, "2024-01-01T11:00:00Z")
	f := filter.New(p, start, end)

	line := "no timestamp here at all"
	if got := f.Evaluate(line); got != filter.Unparseable {
		t.Errorf("expected Unparseable, got %v", got)
	}
}

func TestInRange(t *testing.T) {
	p := newParser(t)
	start := mustTime(t, "2024-01-01T10:00:00Z")
	end := mustTime(t, "2024-01-01T11:00:00Z")
	f := filter.New(p, start, end)

	if !f.InRange("2024-01-01T10:00:00Z INFO boundary start") {
		t.Error("expected start boundary to be in range")
	}
	if !f.InRange("2024-01-01T11:00:00Z INFO boundary end") {
		t.Error("expected end boundary to be in range")
	}
	if f.InRange("2024-01-01T09:59:59Z INFO just before") {
		t.Error("expected line before range to be out of range")
	}
}
