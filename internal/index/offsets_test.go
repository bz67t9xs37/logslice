package index

import (
	"testing"
	"time"
)

func makeEntries(timestamps []string) Entries {
	entries := make(Entries, len(timestamps))
	for i, ts := range timestamps {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			panic(err)
		}
		entries[i] = Entry{Timestamp: t, Offset: int64(i * 100)}
	}
	return entries
}

func TestFindRange_FullCoverage(t *testing.T) {
	e := makeEntries([]string{
		"2024-01-01T10:00:00Z",
		"2024-01-01T11:00:00Z",
		"2024-01-01T12:00:00Z",
	})
	from, _ := time.Parse(time.RFC3339, "2024-01-01T09:00:00Z")
	to, _ := time.Parse(time.RFC3339, "2024-01-01T13:00:00Z")

	start, end, ok := e.FindRange(from, to)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if start != 0 {
		t.Errorf("start: got %d, want 0", start)
	}
	if end != -1 {
		t.Errorf("end: got %d, want -1 (EOF)", end)
	}
}

func TestFindRange_SubRange(t *testing.T) {
	e := makeEntries([]string{
		"2024-01-01T10:00:00Z",
		"2024-01-01T11:00:00Z",
		"2024-01-01T12:00:00Z",
		"2024-01-01T13:00:00Z",
	})
	from, _ := time.Parse(time.RFC3339, "2024-01-01T11:00:00Z")
	to, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")

	start, end, ok := e.FindRange(from, to)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if start != 100 {
		t.Errorf("start: got %d, want 100", start)
	}
	if end != 300 {
		t.Errorf("end: got %d, want 300", end)
	}
}

func TestFindRange_NoMatch(t *testing.T) {
	e := makeEntries([]string{
		"2024-01-01T10:00:00Z",
		"2024-01-01T11:00:00Z",
	})
	from, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	to, _ := time.Parse(time.RFC3339, "2024-01-01T13:00:00Z")

	_, _, ok := e.FindRange(from, to)
	if ok {
		t.Fatal("expected ok=false for out-of-range query")
	}
}

func TestFindRange_EmptyEntries(t *testing.T) {
	var e Entries
	from, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
	to, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")

	_, _, ok := e.FindRange(from, to)
	if ok {
		t.Fatal("expected ok=false for empty entries")
	}
}
