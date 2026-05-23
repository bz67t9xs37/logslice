package index

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeSummaryEntries() []Entry {
	base := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: base, Offset: 0},
		{Timestamp: base.Add(10 * time.Minute), Offset: 100},
		{Timestamp: base.Add(20 * time.Minute), Offset: 200},
		{Timestamp: base.Add(30 * time.Minute), Offset: 300},
	}
}

func TestSummarize_BasicStats(t *testing.T) {
	entries := makeSummaryEntries()
	s := Summarize(entries, 2)

	if s.FileCount != 2 {
		t.Errorf("FileCount: got %d, want 2", s.FileCount)
	}
	if s.EntryCount != 4 {
		t.Errorf("EntryCount: got %d, want 4", s.EntryCount)
	}
	if s.Span != 30*time.Minute {
		t.Errorf("Span: got %s, want 30m", s.Span)
	}
	expectedDensity := 4.0 / 30.0
	if diff := s.Density - expectedDensity; diff > 0.001 || diff < -0.001 {
		t.Errorf("Density: got %.4f, want %.4f", s.Density, expectedDensity)
	}
}

func TestSummarize_EmptyEntries(t *testing.T) {
	s := Summarize(nil, 1)
	if s.EntryCount != 0 {
		t.Errorf("expected 0 entries, got %d", s.EntryCount)
	}
	if s.FileCount != 1 {
		t.Errorf("expected FileCount 1, got %d", s.FileCount)
	}
	if !s.EarliestTime.IsZero() {
		t.Errorf("expected zero EarliestTime")
	}
}

func TestSummarize_SingleEntry(t *testing.T) {
	t0 := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	s := Summarize([]Entry{{Timestamp: t0, Offset: 0}}, 1)

	if s.Span != 0 {
		t.Errorf("Span: got %s, want 0", s.Span)
	}
	if s.Density != 0 {
		t.Errorf("Density: got %.4f, want 0 (zero span)", s.Density)
	}
}

func TestSummary_Fprint(t *testing.T) {
	entries := makeSummaryEntries()
	s := Summarize(entries, 1)

	var buf bytes.Buffer
	s.Fprint(&buf)
	out := buf.String()

	for _, want := range []string{"Files", "Entries", "Earliest", "Latest", "Span", "Density"} {
		if !strings.Contains(out, want) {
			t.Errorf("Fprint output missing %q", want)
		}
	}
}

func TestSummary_Fprint_Empty(t *testing.T) {
	s := Summarize(nil, 0)
	var buf bytes.Buffer
	s.Fprint(&buf)
	if !strings.Contains(buf.String(), "No entries") {
		t.Errorf("expected 'No entries' in empty summary output")
	}
}
