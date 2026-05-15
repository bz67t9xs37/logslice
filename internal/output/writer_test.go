package output_test

import (
	"bytes"
	"strings"
	"testing"

	"logslice/internal/output"
)

func TestWriteLine_SingleLine(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	if err := w.WriteLine("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := buf.String()
	if got != "hello world\n" {
		t.Errorf("expected %q, got %q", "hello world\n", got)
	}
}

func TestWriteLine_MultipleLines(t *testing.T) {
	lines := []string{
		"2024-01-01T00:00:00Z INFO starting service",
		"2024-01-01T00:00:01Z DEBUG config loaded",
		"2024-01-01T00:00:02Z ERROR connection refused",
	}

	var buf bytes.Buffer
	w := output.New(&buf)

	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("unexpected error writing line: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
	for i, want := range lines {
		if got[i] != want {
			t.Errorf("line %d: expected %q, got %q", i, want, got[i])
		}
	}
}

func TestLinesWritten(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	for i := 0; i < 5; i++ {
		_ = w.WriteLine("line")
	}

	if w.LinesWritten() != 5 {
		t.Errorf("expected 5 lines written, got %d", w.LinesWritten())
	}
}

func TestClose_FlushesBuffer(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	_ = w.WriteLine("buffered line")
	if err := w.Close(); err != nil {
		t.Fatalf("close error: %v", err)
	}

	if !strings.Contains(buf.String(), "buffered line") {
		t.Error("expected buffered line to be flushed on Close")
	}
}
