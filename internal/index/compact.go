package index

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// CompactOptions controls how compaction behaves.
type CompactOptions struct {
	// BucketSize is the time window used to bucket entries.
	// Entries within the same bucket are merged into a single representative entry.
	BucketSize time.Duration
	// MaxEntries caps the total number of entries after compaction.
	// Zero means no cap.
	MaxEntries int
}

// CompactResult describes what happened during compaction.
type CompactResult struct {
	InputCount  int
	OutputCount int
	Buckets     int
}

// Fprint writes a human-readable summary of the result to w.
func (r CompactResult) Fprint(w io.Writer) {
	fmt.Fprintf(w, "compaction: %d → %d entries across %d buckets\n",
		r.InputCount, r.OutputCount, r.Buckets)
}

// Compact reduces the number of index entries by bucketing them into
// fixed-size time windows. Within each bucket the first offset is kept
// so that seeking into the bucket always lands before the desired time.
func Compact(entries []Entry, opts CompactOptions) ([]Entry, CompactResult) {
	if len(entries) == 0 || opts.BucketSize <= 0 {
		return entries, CompactResult{InputCount: len(entries), OutputCount: len(entries)}
	}

	// Sort defensively; entries should already be ordered.
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	buckets := map[int64]Entry{}
	for _, e := range sorted {
		key := e.Timestamp.Truncate(opts.BucketSize).UnixNano()
		if existing, ok := buckets[key]; !ok {
			buckets[key] = e
		} else if e.Offset < existing.Offset {
			// Keep the earliest offset within the bucket.
			buckets[key] = e
		}
	}

	out := make([]Entry, 0, len(buckets))
	for _, e := range buckets {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp.Before(out[j].Timestamp)
	})

	if opts.MaxEntries > 0 && len(out) > opts.MaxEntries {
		out = out[:opts.MaxEntries]
	}

	return out, CompactResult{
		InputCount:  len(entries),
		OutputCount: len(out),
		Buckets:     len(buckets),
	}
}
