package index

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeStats(n int, first, last time.Time) []Entry {
	entries := make([]Entry, n)
	for i := 0; i < n; i++ {
		frac := float64(i) / float64(n-1)
		delta := time.Duration(float64(last.Sub(first)) * frac)
		entries[i] = Entry{Timestamp: first.Add(delta), Offset: int64(i * 100)}
	}
	return entries
}

func TestCollect_ReturnsCorrectStats(t *testing.T) {
	first := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	entries := makeStats(10, first, last)

	s, err := Collect(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Entries != 10 {
		t.Errorf("Entries: got %d, want 10", s.Entries)
	}
	if !s.First.Equal(first) {
		t.Errorf("First: got %v, want %v", s.First, first)
	}
	if !s.Last.Equal(last) {
		t.Errorf("Last: got %v, want %v", s.Last, last)
	}
	if s.Span != time.Hour {
		t.Errorf("Span: got %v, want 1h", s.Span)
	}
}

func TestCollect_EmptyEntries(t *testing.T) {
	_, err := Collect(nil)
	if err == nil {
		t.Fatal("expected error for empty entries, got nil")
	}
}

func TestCollect_SingleEntry(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	entries := []Entry{{Timestamp: now, Offset: 0}}
	s, err := Collect(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Span != 0 {
		t.Errorf("Span: got %v, want 0", s.Span)
	}
	if s.Entries != 1 {
		t.Errorf("Entries: got %d, want 1", s.Entries)
	}
}

func TestStats_Fprint(t *testing.T) {
	first := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	last := time.Date(2024, 6, 1, 13, 30, 0, 0, time.UTC)
	entries := makeStats(5, first, last)
	s, _ := Collect(entries)

	var buf bytes.Buffer
	s.Fprint(&buf)
	out := buf.String()

	for _, want := range []string{"entries", "first", "last", "span", "1h30m"} {
		if !strings.Contains(out, want) {
			t.Errorf("Fprint output missing %q; got:\n%s", want, out)
		}
	}
}
