// Package index builds and caches a time-offset index for structured log
// files, enabling fast seek-based extraction of time-range windows without
// scanning the entire file.
//
// Build walks a log file and records the byte offset and parsed timestamp for
// each line, returning a slice of Entry values sorted by time. The index can
// be persisted to disk with SaveCache and reloaded with LoadCache; staleness
// is detected via the source file's modification time.
//
// FindRange uses binary search over an Entry slice to locate the half-open
// interval [start, end) so callers can seek directly to the relevant portion
// of the file.
package index
