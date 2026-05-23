// Package index provides functionality for building, caching, and querying
// byte-offset indexes over structured log files.
//
// # Validation
//
// The Validate function inspects a slice of [Entry] values and reports
// consistency issues that may indicate a corrupt or incorrectly built index:
//
//   - Zero timestamps: entries whose timestamp is the zero value of time.Time
//     are flagged as invalid and cannot be used for range queries.
//
//   - Negative offsets: byte offsets must be non-negative file positions.
//
//   - Out-of-order entries: entries must be monotonically non-decreasing by
//     timestamp; a reversal indicates a log file that was not appended in
//     chronological order or a merge error.
//
//   - Duplicate offsets: two entries pointing to the same byte offset are
//     almost certainly a bug in the indexing or merge pipeline.
//
// Validate returns a [ValidationResult] summarising all issues found, and a
// non-nil error when at least one issue is present.
package index
