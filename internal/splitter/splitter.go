// Package splitter wires together the scanner, filter, and output components
// to extract a time-range window from a log file.
package splitter

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"logslice/internal/filter"
	"logslice/internal/index"
	"logslice/internal/output"
	"logslice/internal/timeparse"
)

// Config holds the parameters for a single split operation.
type Config struct {
	InputPath  string
	OutputPath string
	Start      string
	End        string
	// SampleEvery controls the index sampling interval in bytes (0 = default 1 MiB).
	SampleEvery int64
}

// Run opens the input file, builds a sparse index for fast seeking, then
// streams matching lines to the output file.
func Run(cfg Config) error {
	parser := timeparse.New()

	start, end, err := parser.ParseRange(cfg.Start, cfg.End)
	if err != nil {
		return fmt.Errorf("invalid time range: %w", err)
	}

	f, err := os.Open(cfg.InputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat input: %w", err)
	}

	idx, err := index.Build(f, info.Size(), cfg.SampleEvery, func(line string) (time.Time, error) {
		return parser.Parse(line)
	})
	if err != nil {
		return fmt.Errorf("build index: %w", err)
	}

	seekOffset := idx.FindOffset(start)
	if _, err := f.Seek(seekOffset, 0); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	w, err := output.NewFile(cfg.OutputPath)
	if err != nil {
		return fmt.Errorf("open output: %w", err)
	}
	defer w.Close()

	fl := filter.New(parser, start, end)
	scanner := bufio.NewScanner(f)
	pastWindow := false

	for scanner.Scan() {
		line := scanner.Text()
		result := fl.Evaluate(line)
		switch result {
		case filter.Inside:
			if err := w.WriteLine(line); err != nil {
				return fmt.Errorf("write: %w", err)
			}
		case filter.After:
			pastWindow = true
		}
		if pastWindow {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	return nil
}
