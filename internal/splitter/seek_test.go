package splitter

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"logslice/internal/index"
	"logslice/internal/output"
)

func writeTempLogFile(t *testing.T, lines []string) (path string, entries index.Entries) {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var offset int64
	for _, l := range lines {
		ts, _, _ := strings.Cut(l, " ")
		t2, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			t.Fatalf("bad timestamp in test line %q: %v", l, err)
		}
		entries = append(entries, index.Entry{Timestamp: t2, Offset: offset})
		n, _ := f.WriteString(l + "\n")
		offset += int64(n)
	}
	return f.Name(), entries
}

func TestRunWithIndex_SubRange(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z level=info msg=boot",
		"2024-01-01T11:00:00Z level=info msg=ready",
		"2024-01-01T12:00:00Z level=warn msg=slow",
		"2024-01-01T13:00:00Z level=info msg=done",
	}
	path, entries := writeTempLogFile(t, lines)

	var buf bytes.Buffer
	w := output.New(&buf)

	err := RunWithIndex(path, entries, "2024-01-01T11:00:00Z", "2024-01-01T12:00:00Z", w)
	if err != nil {
		t.Fatalf("RunWithIndex: %v", err)
	}
	w.Close()

	got := buf.String()
	if !strings.Contains(got, "msg=ready") {
		t.Errorf("expected msg=ready in output, got: %q", got)
	}
	if !strings.Contains(got, "msg=slow") {
		t.Errorf("expected msg=slow in output, got: %q", got)
	}
	if strings.Contains(got, "msg=boot") {
		t.Errorf("unexpected msg=boot in output, got: %q", got)
	}
	if strings.Contains(got, "msg=done") {
		t.Errorf("unexpected msg=done in output, got: %q", got)
	}
}

func TestRunWithIndex_NoMatch(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z level=info msg=boot",
	}
	path, entries := writeTempLogFile(t, lines)

	var buf bytes.Buffer
	w := output.New(&buf)

	err := RunWithIndex(path, entries, "2024-01-01T12:00:00Z", "2024-01-01T13:00:00Z", w)
	if err != nil {
		t.Fatalf("RunWithIndex: %v", err)
	}
	w.Close()

	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %q", buf.String())
	}
}
