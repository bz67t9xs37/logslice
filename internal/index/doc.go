// Package index builds and queries a time-ordered index of byte offsets within
// structured log files. Each index entry maps a parsed timestamp to the byte
// offset of the corresponding log line, enabling fast seek-based extraction of
// time-range windows without scanning the entire file.
//
// # Building
//
// Use [Build] to construct an index from an io.ReadSeeker. The resulting
// []Entry slice is sorted by timestamp. Indexes can be persisted with
// [SaveCache] and retrieved with [LoadCache]; the cache is invalidated
// automatically when the source file is modified.
//
// # Querying
//
// [FindRange] returns the start and end byte offsets for a given time window.
// [Union] merges indexes from multiple files into a single sorted slice.
// [MergeEntries] combines two pre-sorted slices, deduplicating by offset.
// [Trim] restricts an index to a specific time window.
// [Prune] removes stale or excess entries for cache-management purposes.
//
// # Reporting
//
// [Collect] computes aggregate [Stats] (min, max, count, span) from an index.
// [Summarize] produces a [Summary] with human-readable density metrics.
