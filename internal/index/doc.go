// Package index builds and queries a sparse timestamp index over a structured
// log file, enabling O(log n) seek-to-time without scanning the entire file.
//
// # Building an index
//
// Call [Build] with the path to a log file and a timestamp-parsing function.
// Build scans every line, extracts its timestamp, and records the byte offset
// of each line that successfully parses. The resulting []Entry slice is sorted
// by timestamp and can be queried with [FindOffset].
//
// # Persistent caching
//
// Repeated runs over the same large file can be expensive. The cache sub-API
// ([LoadCache], [SaveCache], [CachePath]) persists the index alongside the
// source file as a hidden gob-encoded file. The cache is validated against the
// file's mtime and size; a mismatch causes a transparent cache miss so callers
// always receive a correct index.
package index
