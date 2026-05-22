package index

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry holds a persisted index alongside metadata used to
// validate whether the cache is still fresh.
type CacheEntry struct {
	ModTime time.Time
	Size    int64
	Entries []Entry
}

// CachePath returns the conventional cache file path for the given log file.
// The cache is stored as a hidden file alongside the source with a .idx suffix.
func CachePath(logPath string) string {
	dir := filepath.Dir(logPath)
	base := filepath.Base(logPath)
	return filepath.Join(dir, "."+base+".idx")
}

// LoadCache reads a previously persisted index from disk and validates it
// against the current stat of logPath. It returns nil, nil when no valid
// cache exists so callers can transparently fall back to Build.
func LoadCache(logPath string) ([]Entry, error) {
	info, err := os.Stat(logPath)
	if err != nil {
		return nil, fmt.Errorf("stat log file: %w", err)
	}

	f, err := os.Open(CachePath(logPath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open cache: %w", err)
	}
	defer f.Close()

	var ce CacheEntry
	if err := gob.NewDecoder(f).Decode(&ce); err != nil {
		return nil, nil // treat corrupt cache as a miss
	}

	if ce.ModTime != info.ModTime() || ce.Size != info.Size() {
		return nil, nil // stale
	}

	return ce.Entries, nil
}

// SaveCache writes entries to the conventional cache path for logPath.
func SaveCache(logPath string, entries []Entry) error {
	info, err := os.Stat(logPath)
	if err != nil {
		return fmt.Errorf("stat log file: %w", err)
	}

	ce := CacheEntry{
		ModTime: info.ModTime(),
		Size:    info.Size(),
		Entries: entries,
	}

	tmp := CachePath(logPath) + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("create cache tmp: %w", err)
	}

	if err := gob.NewEncoder(f).Encode(ce); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("encode cache: %w", err)
	}
	f.Close()

	if err := os.Rename(tmp, CachePath(logPath)); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("rename cache: %w", err)
	}
	return nil
}
