package index

import (
	"bytes"
	"testing"
	"time"
)

func makeCompactEntries(base time.Time, count int, step time.Duration) []Entry {
	ent := make([]Entry, count)
	for i := range ent {
		ent[i] = Entry{
			Timestamp: base.Add(time.Duration(i) * step),
			Offset:    int64(i * 100),
		}
	}
	return ent
}

func TestCompact_ReducesByBucket(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// 60 entries, one per minute → 6 buckets of 10 min each
	entries := makeCompactEntries(base, 60, time.Minute)
	out, res := Compact(entries, CompactOptions{BucketSize: 10 * time.Minute})
	if res.InputCount != 60 {
		t.Errorf("InputCount: got %d, want 60", res.InputCount)
	}
	if len(out) != 6 {
		t.Errorf("output len: got %d, want 6", len(out))
	}
	if res.Buckets != 6 {
		t.Errorf("Buckets: got %d, want 6", res.Buckets)
	}
}

func TestCompact_MaxEntries(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := makeCompactEntries(base, 60, time.Minute)
	out, res := Compact(entries, CompactOptions{BucketSize: time.Minute, MaxEntries: 10})
	if len(out) != 10 {
		t.Errorf("output len: got %d, want 10", len(out))
	}
	if res.OutputCount != 10 {
		t.Errorf("OutputCount: got %d, want 10", res.OutputCount)
	}
}

func TestCompact_EmptyEntries(t *testing.T) {
	out, res := Compact(nil, CompactOptions{BucketSize: time.Minute})
	if len(out) != 0 {
		t.Errorf("expected empty output")
	}
	if res.InputCount != 0 || res.OutputCount != 0 {
		t.Errorf("unexpected counts: %+v", res)
	}
}

func TestCompact_ZeroBucketSize_NoOp(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := makeCompactEntries(base, 5, time.Minute)
	out, res := Compact(entries, CompactOptions{BucketSize: 0})
	if len(out) != 5 {
		t.Errorf("expected passthrough, got %d entries", len(out))
	}
	if res.InputCount != res.OutputCount {
		t.Errorf("counts should match for no-op")
	}
}

func TestCompact_KeepsEarliestOffset(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: base.Add(1 * time.Minute), Offset: 200},
		{Timestamp: base.Add(2 * time.Minute), Offset: 50},
		{Timestamp: base.Add(3 * time.Minute), Offset: 300},
	}
	out, _ := Compact(entries, CompactOptions{BucketSize: 10 * time.Minute})
	if len(out) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(out))
	}
	if out[0].Offset != 50 {
		t.Errorf("expected earliest offset 50, got %d", out[0].Offset)
	}
}

func TestCompactResult_Fprint(t *testing.T) {
	r := CompactResult{InputCount: 100, OutputCount: 10, Buckets: 10}
	var buf bytes.Buffer
	r.Fprint(&buf)
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
