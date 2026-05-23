package index

import (
	"bytes"
	"testing"
	"time"
)

func makeRepairEntries() []Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: base.Add(2 * time.Second), Offset: 200},
		{Timestamp: base.Add(1 * time.Second), Offset: 100},
		{Timestamp: base.Add(3 * time.Second), Offset: 300},
	}
}

func TestRepair_AlreadyClean(t *testing.T) {
	entries := makeRepairEntries()
	// Sort them first so no reorder needed.
	entries[0], entries[1] = entries[1], entries[0]
	out, res := Repair(entries)
	if res.RemovedZero != 0 || res.RemovedDups != 0 || res.Reordered {
		t.Errorf("expected no changes, got %+v", res)
	}
	if res.RepairedCount != 3 {
		t.Errorf("expected 3 entries, got %d", res.RepairedCount)
	}
	_ = out
}

func TestRepair_RemovesZeroTimestamps(t *testing.T) {
	entries := makeRepairEntries()
	entries = append(entries, Entry{Timestamp: time.Time{}, Offset: 999})
	out, res := Repair(entries)
	if res.RemovedZero != 1 {
		t.Errorf("expected 1 zero removed, got %d", res.RemovedZero)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries after repair, got %d", len(out))
	}
}

func TestRepair_RemovesDuplicateOffsets(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: base.Add(1 * time.Second), Offset: 100},
		{Timestamp: base.Add(2 * time.Second), Offset: 100}, // duplicate offset
		{Timestamp: base.Add(3 * time.Second), Offset: 200},
	}
	out, res := Repair(entries)
	if res.RemovedDups != 1 {
		t.Errorf("expected 1 dup removed, got %d", res.RemovedDups)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestRepair_SortsOutOfOrder(t *testing.T) {
	entries := makeRepairEntries() // out of order
	out, res := Repair(entries)
	if !res.Reordered {
		t.Error("expected Reordered=true")
	}
	for i := 1; i < len(out); i++ {
		if out[i].Timestamp.Before(out[i-1].Timestamp) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestRepairResult_Fprint(t *testing.T) {
	r := RepairResult{OriginalCount: 10, RepairedCount: 8, RemovedZero: 1, RemovedDups: 1, Reordered: true}
	var buf bytes.Buffer
	r.Fprint(&buf)
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
