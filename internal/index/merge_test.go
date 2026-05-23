package index

import (
	"testing"
	"time"
)

func ts(sec int) time.Time {
	return time.Unix(int64(sec), 0).UTC()
}

func TestMergeEntries_Disjoint(t *testing.T) {
	a := []Entry{{Timestamp: ts(1), Offset: 0}, {Timestamp: ts(3), Offset: 20}}
	b := []Entry{{Timestamp: ts(2), Offset: 10}, {Timestamp: ts(4), Offset: 30}}

	got := MergeEntries(a, b)
	if len(got) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(got))
	}
	for i := 1; i < len(got); i++ {
		if got[i].Timestamp.Before(got[i-1].Timestamp) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestMergeEntries_Deduplication(t *testing.T) {
	a := []Entry{{Timestamp: ts(1), Offset: 0}, {Timestamp: ts(2), Offset: 10}}
	b := []Entry{{Timestamp: ts(2), Offset: 10}, {Timestamp: ts(3), Offset: 20}}

	got := MergeEntries(a, b)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries after dedup, got %d", len(got))
	}
}

func TestMergeEntries_EmptyA(t *testing.T) {
	b := []Entry{{Timestamp: ts(1), Offset: 0}}
	got := MergeEntries(nil, b)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}

func TestMergeEntries_EmptyB(t *testing.T) {
	a := []Entry{{Timestamp: ts(1), Offset: 0}}
	got := MergeEntries(a, nil)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}

func TestTrim_SubRange(t *testing.T) {
	entries := []Entry{
		{Timestamp: ts(1), Offset: 0},
		{Timestamp: ts(5), Offset: 40},
		{Timestamp: ts(10), Offset: 80},
		{Timestamp: ts(15), Offset: 120},
	}
	got := Trim(entries, ts(5), ts(10))
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Offset != 40 || got[1].Offset != 80 {
		t.Errorf("unexpected offsets: %v", got)
	}
}

func TestTrim_NoMatch(t *testing.T) {
	entries := []Entry{
		{Timestamp: ts(1), Offset: 0},
		{Timestamp: ts(2), Offset: 10},
	}
	got := Trim(entries, ts(5), ts(10))
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestTrim_AllMatch(t *testing.T) {
	entries := []Entry{
		{Timestamp: ts(3), Offset: 0},
		{Timestamp: ts(7), Offset: 50},
	}
	got := Trim(entries, ts(1), ts(10))
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}
