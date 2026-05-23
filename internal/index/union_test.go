package index

import (
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempLogForUnion(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "union-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, _ = f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func TestUnion_SingleFile(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z info startup",
		"2024-01-01T10:01:00Z info ready",
		"2024-01-01T10:02:00Z info request",
	}
	path := writeTempLogForUnion(t, lines)

	start := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 10, 2, 0, 0, time.UTC)

	res, err := Union([]string{path}, start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(res.Sources))
	}
	if len(res.Entries) == 0 {
		t.Error("expected non-empty entries")
	}
}

func TestUnion_MultipleFiles_Merged(t *testing.T) {
	pathA := writeTempLogForUnion(t, []string{
		"2024-01-01T10:00:00Z info a1",
		"2024-01-01T10:02:00Z info a2",
	})
	pathB := writeTempLogForUnion(t, []string{
		"2024-01-01T10:01:00Z info b1",
		"2024-01-01T10:03:00Z info b2",
	})

	start := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 10, 3, 0, 0, time.UTC)

	res, err := Union([]string{pathA, pathB}, start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(res.Sources))
	}
	for i := 1; i < len(res.Entries); i++ {
		if res.Entries[i].Timestamp.Before(res.Entries[i-1].Timestamp) {
			t.Errorf("merged entries not sorted at index %d", i)
		}
	}
}

func TestUnion_NoPaths(t *testing.T) {
	_, err := Union(nil, time.Now(), time.Now())
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestUnion_NoMatchingEntries(t *testing.T) {
	path := writeTempLogForUnion(t, []string{
		"2024-01-01T08:00:00Z info old",
	})
	start := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)

	res, err := Union([]string{path}, start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(res.Entries))
	}
	if len(res.Sources) != 0 {
		t.Errorf("expected 0 sources, got %d", len(res.Sources))
	}
}
