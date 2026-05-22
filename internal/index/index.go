// Package index provides byte-offset indexing for log files,
// enabling fast seeks to approximate time positions without
// scanning the entire file from the start.
package index

import (
	"io"
	"time"
)

// Entry records the byte offset of a log line along with its parsed timestamp.
type Entry struct {
	Offset    int64
	Timestamp time.Time
}

// Index is an ordered collection of sampled offset entries for a log file.
type Index struct {
	entries []Entry
}

// Build reads from r, sampling one entry every sampleEvery bytes, using
// parseFn to extract timestamps from each line. Lines that cannot be parsed
// are skipped. r must support ReadAt (e.g. *os.File).
func Build(r io.ReaderAt, size int64, sampleEvery int64, parseFn func(string) (time.Time, error)) (*Index, error) {
	if sampleEvery <= 0 {
		sampleEvery = 1 << 20 // 1 MiB default
	}

	idx := &Index{}
	var offset int64

	for offset < size {
		line, newOffset, err := readLineAt(r, offset, size)
		if err != nil {
			break
		}
		if ts, err := parseFn(line); err == nil {
			idx.entries = append(idx.entries, Entry{Offset: offset, Timestamp: ts})
		}
		nextSample := offset + sampleEvery
		offset, _, err = seekToLineStart(r, nextSample, size)
		if err != nil || offset >= size {
			break
		}
		_ = newOffset
	}
	return idx, nil
}

// FindOffset returns the byte offset of the last index entry whose timestamp
// is <= target, or 0 if no suitable entry exists.
func (idx *Index) FindOffset(target time.Time) int64 {
	var best int64
	for _, e := range idx.entries {
		if !e.Timestamp.After(target) {
			best = e.Offset
		} else {
			break
		}
	}
	return best
}

// Len returns the number of entries in the index.
func (idx *Index) Len() int { return len(idx.entries) }

// readLineAt reads a single line starting at offset from r.
func readLineAt(r io.ReaderAt, offset, size int64) (string, int64, error) {
	const bufSize = 4096
	buf := make([]byte, bufSize)
	n, err := r.ReadAt(buf, offset)
	if n == 0 {
		return "", offset, err
	}
	for i := 0; i < n; i++ {
		if buf[i] == '\n' {
			return string(buf[:i]), offset + int64(i) + 1, nil
		}
	}
	return string(buf[:n]), offset + int64(n), nil
}

// seekToLineStart advances to the start of the next complete line at or after pos.
func seekToLineStart(r io.ReaderAt, pos, size int64) (int64, bool, error) {
	if pos >= size {
		return pos, false, io.EOF
	}
	const bufSize = 512
	buf := make([]byte, bufSize)
	n, err := r.ReadAt(buf, pos)
	for i := 0; i < n; i++ {
		if buf[i] == '\n' {
			return pos + int64(i) + 1, true, nil
		}
	}
	if err != nil {
		return pos, false, err
	}
	return pos + int64(n), false, nil
}
