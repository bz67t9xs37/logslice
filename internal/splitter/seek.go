// Package splitter provides the top-level log splitting pipeline.
package splitter

import (
	"fmt"
	"io"
	"os"

	"logslice/internal/index"
	"logslice/internal/output"
	"logslice/internal/scanner"
	"logslice/internal/timeparse"
)

// RunWithIndex is like Run but uses a pre-built index to seek directly to the
// relevant byte range, skipping lines outside [start, end].
// This can dramatically reduce I/O for large log files.
func RunWithIndex(
	logPath string,
	entries index.Entries,
	startStr, endStr string,
	w *output.Writer,
) error {
	p := timeparse.New()
	tr, err := p.ParseRange(startStr, endStr)
	if err != nil {
		return fmt.Errorf("parse range: %w", err)
	}

	startOff, endOff, ok := entries.FindRange(tr.Start, tr.End)
	if !ok {
		// No entries in range — nothing to write.
		return nil
	}

	f, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("open log: %w", err)
	}
	defer f.Close()

	if _, err := f.Seek(startOff, io.SeekStart); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	var r io.Reader = f
	if endOff >= 0 {
		r = io.LimitReader(f, endOff-startOff)
	}

	sc := scanner.New(r, p)
	for sc.Scan() {
		line := sc.Line()
		if err := w.WriteLine(line.Raw); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}
	return sc.Err()
}
