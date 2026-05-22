// Package index implements sparse byte-offset indexing for structured log files.
//
// Building an index samples the log file at regular byte intervals, recording
// the file offset and parsed timestamp of each sampled line. The index can then
// be queried to find the best byte offset to seek to before beginning a
// time-range scan, avoiding a full sequential read from the start of the file.
//
// Typical usage:
//
//	idx, err := index.Build(f, size, 1<<20, parser.Parse)
//	if err != nil { ... }
//	seekOffset := idx.FindOffset(startTime)
//	// seek f to seekOffset, then hand off to scanner
//
// The index is held entirely in memory and is not persisted; it is intended to
// be rebuilt cheaply on each invocation for files that fit the sampling budget.
package index
