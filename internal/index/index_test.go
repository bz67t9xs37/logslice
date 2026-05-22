package index

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// readerAt wraps a string so it satisfies io.ReaderAt.
type readerAt struct{ s string }

func (r *readerAt) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(r.s)) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, r.s[off:])
	return n, nil
}

func makeLogData(start time.Time, count int, interval time.Duration) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		ts := start.Add(time.Duration(i) * interval)
		fmt.Fprintf(&sb, "%s log line %d\n", ts.Format(time.RFC3339), i)
	}
	return sb.String()
}

func parseTS(line string) (time.Time, error) {
	if len(line) < 20 {
		return time.Time{}, fmt.Errorf("too short")
	}
	return time.Parse(time.RFC3339, line[:20])
}

func TestBuild_ReturnsEntries(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	data := makeLogData(base, 50, time.Minute)
	r := &readerAt{s: data}

	idx, err := Build(r, int64(len(data)), 200, parseTS)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if idx.Len() == 0 {
		t.Fatal("expected at least one index entry")
	}
}

func TestFindOffset_BeforeAll(t *testing.T) {
	base := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	data := makeLogData(base, 20, time.Minute)
	r := &readerAt{s: data}

	idx, _ := Build(r, int64(len(data)), 100, parseTS)
	off := idx.FindOffset(base.Add(-time.Hour))
	if off != 0 {
		t.Errorf("expected offset 0 for timestamp before all entries, got %d", off)
	}
}

func TestFindOffset_AfterAll(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	data := makeLogData(base, 20, time.Minute)
	r := &readerAt{s: data}

	idx, _ := Build(r, int64(len(data)), 100, parseTS)
	prev := idx.FindOffset(base.Add(24 * time.Hour))
	if prev == 0 && idx.Len() > 1 {
		t.Error("expected non-zero offset for timestamp after all entries")
	}
}

func TestFindOffset_Monotonic(t *testing.T) {
	base := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	data := makeLogData(base, 100, 30*time.Second)
	r := &readerAt{s: data}

	idx, _ := Build(r, int64(len(data)), 150, parseTS)

	var lastOff int64
	for i := 0; i < 100; i++ {
		t2 := base.Add(time.Duration(i) * 30 * time.Second)
		off := idx.FindOffset(t2)
		if off < lastOff {
			t.Errorf("offset not monotonic at step %d: %d < %d", i, off, lastOff)
		}
		lastOff = off
	}
}

func TestLen_EmptyIndex(t *testing.T) {
	idx := &Index{}
	if idx.Len() != 0 {
		t.Errorf("expected 0, got %d", idx.Len())
	}
}
