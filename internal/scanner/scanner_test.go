package scanner_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/scanner"
	"github.com/logslice/logslice/internal/timeparse"
)

const sampleLog = `2024-01-15T10:00:00Z INFO  service started
2024-01-15T10:01:00Z DEBUG request received id=1
2024-01-15T10:02:00Z INFO  processed id=1
2024-01-15T10:03:00Z WARN  slow query detected
2024-01-15T10:04:00Z ERROR connection refused
2024-01-15T10:05:00Z INFO  retry succeeded
`

func newParser(t *testing.T) *timeparse.Parser {
	t.Helper()
	p, err := timeparse.New(nil)
	if err != nil {
		t.Fatalf("timeparse.New: %v", err)
	}
	return p
}

func collectLines(t *testing.T, log string, start, end time.Time) []string {
	t.Helper()
	p := newParser(t)
	s := scanner.New(p, start, end)

	var lines []string
	err := s.Scan(strings.NewReader(log), func(line string) error {
		lines = append(lines, line)
		return nil
	})
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	return lines
}

func TestScan_FullRange(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 10, 5, 0, 0, time.UTC)
	lines := collectLines(t, sampleLog, start, end)
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(lines))
	}
}

func TestScan_SubRange(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 2, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 10, 3, 0, 0, time.UTC)
	lines := collectLines(t, sampleLog, start, end)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
}

func TestScan_NoMatch(t *testing.T) {
	start := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	lines := collectLines(t, sampleLog, start, end)
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestScan_ContinuationLines(t *testing.T) {
	log := "2024-01-15T10:02:00Z INFO start\n  continuation line\n2024-01-15T10:04:00Z INFO end\n"
	start := time.Date(2024, 1, 15, 10, 2, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 10, 3, 0, 0, time.UTC)
	lines := collectLines(t, log, start, end)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (entry + continuation), got %d: %v", len(lines), lines)
	}
}
