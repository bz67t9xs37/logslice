// Package index provides indexing, caching, and querying utilities for
// structured log files used by logslice.
//
// # Compaction
//
// Over time an index can accumulate a large number of entries — one per
// indexed line — which increases memory usage and slows down offset lookups.
// Compaction reduces the entry count by merging entries that fall within the
// same time bucket into a single representative entry whose offset points to
// the earliest position in that bucket.
//
// Usage:
//
//	opts := index.CompactOptions{
//	    BucketSize: 5 * time.Minute,
//	    MaxEntries: 10_000,
//	}
//	compacted, result, err := index.CompactFile("/var/log/app.log", opts)
//	if err != nil { ... }
//	result.Fprint(os.Stdout)
//
// The compacted index is automatically persisted to the cache so that
// subsequent runs benefit from the reduced size without rebuilding.
package index
