package index

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeExportEntries() []Entry {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: base, Offset: 0, Length: 42},
		{Timestamp: base.Add(time.Minute), Offset: 42, Length: 58},
		{Timestamp: base.Add(2 * time.Minute), Offset: 100, Length: 75},
	}
}

func TestExport_JSON(t *testing.T) {
	entries := makeExportEntries()
	var buf bytes.Buffer
	if err := Export(&buf, entries, FormatJSON); err != nil {
		t.Fatalf("Export JSON: %v", err)
	}
	var got []exportEntry
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != len(entries) {
		t.Fatalf("want %d entries, got %d", len(entries), len(got))
	}
	if got[0].Offset != 0 || got[1].Offset != 42 || got[2].Offset != 100 {
		t.Errorf("unexpected offsets: %v", got)
	}
	if got[0].Length != 42 {
		t.Errorf("want length 42, got %d", got[0].Length)
	}
}

func TestExport_CSV(t *testing.T) {
	entries := makeExportEntries()
	var buf bytes.Buffer
	if err := Export(&buf, entries, FormatCSV); err != nil {
		t.Fatalf("Export CSV: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 3 data rows
	if len(lines) != 4 {
		t.Fatalf("want 4 lines, got %d:\n%s", len(lines), buf.String())
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("expected CSV header, got %q", lines[0])
	}
	if !strings.Contains(lines[1], "0") {
		t.Errorf("first data row missing offset 0: %q", lines[1])
	}
}

func TestExport_Text(t *testing.T) {
	entries := makeExportEntries()
	var buf bytes.Buffer
	if err := Export(&buf, entries, FormatText); err != nil {
		t.Fatalf("Export Text: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if !strings.Contains(l, "offset=") || !strings.Contains(l, "length=") {
			t.Errorf("malformed text line: %q", l)
		}
	}
}

func TestExport_EmptyEntries(t *testing.T) {
	for _, fmt := range []ExportFormat{FormatJSON, FormatCSV, FormatText} {
		var buf bytes.Buffer
		if err := Export(&buf, nil, fmt); err != nil {
			t.Errorf("format %s with nil entries: %v", fmt, err)
		}
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, makeExportEntries(), ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
	if !strings.Contains(err.Error(), "unknown format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
