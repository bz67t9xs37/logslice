package index

import (
	"testing"
	"time"
)

func makeEntriesForPrune(timestamps []time.Time) []Entry {
	entries := make([]Entry, len(timestamps))
	for i, t := range timestamps {
		entries[i] = Entry{Timestamp: t, Offset: int64(i * 100)}
	}
	return entries
}

func TestPrune_OlderThan(t *testing.T) {
	now := time.Now()
	times := []time.Time{
		now.Add(-3 * time.Hour),
		now.Add(-2 * time.Hour),
		now.Add(-30 * time.Minute),
		now.Add(-5 * time.Minute),
	}
	entries := makeEntriesForPrune(times)

	result := Prune(entries, PruneOptions{OlderThan: time.Hour})

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if !result[0].Timestamp.Equal(times[2]) {
		t.Errorf("unexpected first entry: %v", result[0].Timestamp)
	}
}

func TestPrune_MaxEntries(t *testing.T) {
	now := time.Now()
	times := []time.Time{
		now.Add(-4 * time.Hour),
		now.Add(-3 * time.Hour),
		now.Add(-2 * time.Hour),
		now.Add(-1 * time.Hour),
		now,
	}
	entries := makeEntriesForPrune(times)

	result := Prune(entries, PruneOptions{MaxEntries: 3})

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if !result[0].Timestamp.Equal(times[2]) {
		t.Errorf("expected most-recent 3 entries; first got %v", result[0].Timestamp)
	}
}

func TestPrune_OlderThanAndMaxEntries(t *testing.T) {
	now := time.Now()
	times := []time.Time{
		now.Add(-5 * time.Hour),
		now.Add(-2 * time.Hour),
		now.Add(-90 * time.Minute),
		now.Add(-20 * time.Minute),
		now.Add(-5 * time.Minute),
	}
	entries := makeEntriesForPrune(times)

	// OlderThan 3h removes first entry; MaxEntries 2 keeps last two.
	result := Prune(entries, PruneOptions{OlderThan: 3 * time.Hour, MaxEntries: 2})

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if !result[0].Timestamp.Equal(times[3]) {
		t.Errorf("unexpected entry: %v", result[0].Timestamp)
	}
}

func TestPrune_EmptyEntries(t *testing.T) {
	result := Prune(nil, PruneOptions{OlderThan: time.Hour, MaxEntries: 10})
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

func TestPrune_NoOptions(t *testing.T) {
	now := time.Now()
	entries := makeEntriesForPrune([]time.Time{now.Add(-1 * time.Hour), now})

	result := Prune(entries, PruneOptions{})

	if len(result) != 2 {
		t.Fatalf("expected all entries preserved, got %d", len(result))
	}
}
