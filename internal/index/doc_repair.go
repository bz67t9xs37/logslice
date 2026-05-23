// Package index provides index building, caching, and management for logslice.
//
// # Repair
//
// The repair sub-feature cleans up a persisted index cache by:
//
//   - Removing entries with zero (uninitialised) timestamps.
//   - Deduplicating entries that share the same byte offset, keeping the
//     entry with the earliest timestamp.
//   - Re-sorting entries in ascending timestamp order if they are found to
//     be out of order.
//
// Repair is non-destructive with respect to the source log file; it only
// modifies the companion cache file produced by [SaveCache].
//
// # Usage
//
//	// High-level helper — repairs and prints a summary.
//	if err := index.RunRepair("/var/log/app.log"); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Lower-level access with options.
//	result, err := index.RepairFile("/var/log/app.log", index.RepairOptions{Verbose: true})
package index
