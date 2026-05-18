// Package splitter orchestrates reading, filtering, and writing log lines
// within a specified time range.
package splitter

import (
	"fmt"
	"io"

	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/scanner"
	"logslice/internal/timeparse"
)

// Config holds the parameters for a split operation.
type Config struct {
	// Input is the source log data to read from.
	Input io.ReadSeeker
	// Output is the destination for matched log lines.
	Output io.Writer
	// Start is the beginning of the time window (inclusive), RFC3339 or similar.
	Start string
	// End is the end of the time window (inclusive).
	End string
}

// Result summarises a completed split operation.
type Result struct {
	LinesScanned int
	LinesWritten int
}

// Run executes the log-splitting pipeline defined by cfg.
func Run(cfg Config) (Result, error) {
	p := timeparse.New()

	start, end, err := p.ParseRange(cfg.Start, cfg.End)
	if err != nil {
		return Result{}, fmt.Errorf("splitter: parse range: %w", err)
	}

	f := filter.New(p, start, end)
	w := output.New(cfg.Output)
	defer w.Close() //nolint:errcheck

	s := scanner.New(cfg.Input, f)

	var res Result
	for {
		line, ok, err := s.Next()
		if err != nil {
			return res, fmt.Errorf("splitter: scan: %w", err)
		}
		if !ok {
			break
		}
		res.LinesScanned++
		if err := w.WriteLine(line); err != nil {
			return res, fmt.Errorf("splitter: write: %w", err)
		}
		res.LinesWritten++
	}

	if err := w.Close(); err != nil {
		return res, fmt.Errorf("splitter: flush: %w", err)
	}
	return res, nil
}
