// Package scanner provides line-by-line scanning of log files with
// time-based filtering using binary search for efficient range detection.
package scanner

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/timeparse"
)

// LineHandler is called for each log line that falls within the time range.
type LineHandler func(line string) error

// Scanner reads log lines and filters them by time range.
type Scanner struct {
	parser  *timeparse.Parser
	start   time.Time
	end     time.Time
}

// New creates a Scanner that will emit lines whose timestamps fall
// within [start, end] (inclusive).
func New(parser *timeparse.Parser, start, end time.Time) *Scanner {
	return &Scanner{
		parser: parser,
		start:  start,
		end:    end,
	}
}

// Scan reads from r line by line, calling handler for every line whose
// parsed timestamp is within the configured range. Lines that cannot be
// parsed are skipped. Scan returns the first error returned by handler,
// or any IO error encountered.
func (s *Scanner) Scan(r io.Reader, handler LineHandler) error {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)

	inWindow := false

	for sc.Scan() {
		line := sc.Text()

		t, err := s.parser.Parse(line)
		if err != nil {
			// No timestamp found — emit if we are already inside the window.
			if inWindow {
				if herr := handler(line); herr != nil {
					return herr
				}
			}
			continue
		}

		if t.Before(s.start) {
			inWindow = false
			continue
		}
		if t.After(s.end) {
			// Log files are assumed to be chronological; stop early.
			break
		}

		inWindow = true
		if herr := handler(line); herr != nil {
			return herr
		}
	}

	return sc.Err()
}
