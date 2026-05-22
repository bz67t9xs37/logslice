package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, data []byte) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "app.log")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatalf("write temp log: %v", err)
	}
	return p
}

func TestCachePath(t *testing.T) {
	got := CachePath("/var/log/app.log")
	want := "/var/log/.app.log.idx"
	if got != want {
		t.Errorf("CachePath = %q; want %q", got, want)
	}
}

func TestSaveAndLoadCache_RoundTrip(t *testing.T) {
	logPath := writeTempLog(t, makeLogData())

	entries, err := Build(logPath, parseTS)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	if err := SaveCache(logPath, entries); err != nil {
		t.Fatalf("SaveCache: %v", err)
	}

	loaded, err := LoadCache(logPath)
	if err != nil {
		t.Fatalf("LoadCache: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected cached entries, got nil")
	}
	if len(loaded) != len(entries) {
		t.Errorf("len(loaded) = %d; want %d", len(loaded), len(entries))
	}
	for i := range entries {
		if loaded[i].Offset != entries[i].Offset {
			t.Errorf("entry[%d].Offset = %d; want %d", i, loaded[i].Offset, entries[i].Offset)
		}
		if !loaded[i].Timestamp.Equal(entries[i].Timestamp) {
			t.Errorf("entry[%d].Timestamp mismatch", i)
		}
	}
}

func TestLoadCache_MissWhenAbsent(t *testing.T) {
	logPath := writeTempLog(t, makeLogData())
	entries, err := LoadCache(logPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil on cache miss, got %d entries", len(entries))
	}
}

func TestLoadCache_StaleAfterModification(t *testing.T) {
	logPath := writeTempLog(t, makeLogData())

	entries, _ := Build(logPath, parseTS)
	_ = SaveCache(logPath, entries)

	// Modify the log file so the cache becomes stale.
	time.Sleep(10 * time.Millisecond)
	f, _ := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	f.WriteString("2024-01-01T00:00:05Z extra line\n")
	f.Close()

	loaded, err := LoadCache(logPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil for stale cache, got entries")
	}
}
