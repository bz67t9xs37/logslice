package splitter_test

import (
	"bytes"
	"strings"
	"testing"

	"logslice/internal/splitter"
)

const sampleLog = `2024-01-10T08:00:00Z INFO  service started
2024-01-10T09:00:00Z DEBUG request received
2024-01-10T10:00:00Z INFO  processing
2024-01-10T11:00:00Z WARN  slow query
2024-01-10T12:00:00Z ERROR timeout
2024-01-10T13:00:00Z INFO  done
`

func TestRun_FullRange(t *testing.T) {
	var out bytes.Buffer
	res, err := splitter.Run(splitter.Config{
		Input:  strings.NewReader(sampleLog),
		Output: &out,
		Start:  "2024-01-10T08:00:00Z",
		End:    "2024-01-10T13:00:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesWritten != 6 {
		t.Errorf("expected 6 lines written, got %d", res.LinesWritten)
	}
}

func TestRun_SubRange(t *testing.T) {
	var out bytes.Buffer
	res, err := splitter.Run(splitter.Config{
		Input:  strings.NewReader(sampleLog),
		Output: &out,
		Start:  "2024-01-10T10:00:00Z",
		End:    "2024-01-10T11:00:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesWritten != 2 {
		t.Errorf("expected 2 lines written, got %d", res.LinesWritten)
	}
	if !strings.Contains(out.String(), "processing") {
		t.Errorf("expected 'processing' in output")
	}
	if !strings.Contains(out.String(), "slow query") {
		t.Errorf("expected 'slow query' in output")
	}
}

func TestRun_NoMatch(t *testing.T) {
	var out bytes.Buffer
	res, err := splitter.Run(splitter.Config{
		Input:  strings.NewReader(sampleLog),
		Output: &out,
		Start:  "2024-01-11T00:00:00Z",
		End:    "2024-01-11T23:59:59Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesWritten != 0 {
		t.Errorf("expected 0 lines written, got %d", res.LinesWritten)
	}
}

func TestRun_InvalidRange(t *testing.T) {
	var out bytes.Buffer
	_, err := splitter.Run(splitter.Config{
		Input:  strings.NewReader(sampleLog),
		Output: &out,
		Start:  "2024-01-10T12:00:00Z",
		End:    "2024-01-10T08:00:00Z",
	})
	if err == nil {
		t.Fatal("expected error for start-after-end range")
	}
}
