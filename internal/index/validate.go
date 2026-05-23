package index

import (
	"errors"
	"fmt"
	"time"
)

// ValidationResult holds the outcome of validating an index entry slice.
type ValidationResult struct {
	TotalEntries  int
	InvalidCount  int
	OutOfOrder    int
	Duplicates    int
	Errors        []string
}

// Valid returns true if no validation issues were found.
func (v ValidationResult) Valid() bool {
	return v.InvalidCount == 0 && v.OutOfOrder == 0 && v.Duplicates == 0
}

// Validate checks a slice of index entries for consistency issues such as
// zero timestamps, out-of-order entries, and duplicate offsets.
func Validate(entries []Entry) (ValidationResult, error) {
	result := ValidationResult{
		TotalEntries: len(entries),
	}

	if len(entries) == 0 {
		return result, nil
	}

	seenOffsets := make(map[int64]bool, len(entries))

	for i, e := range entries {
		if e.Timestamp.IsZero() {
			msg := fmt.Sprintf("entry %d: zero timestamp at offset %d", i, e.Offset)
			result.Errors = append(result.Errors, msg)
			result.InvalidCount++
		}

		if e.Offset < 0 {
			msg := fmt.Sprintf("entry %d: negative offset %d", i, e.Offset)
			result.Errors = append(result.Errors, msg)
			result.InvalidCount++
		}

		if seenOffsets[e.Offset] {
			msg := fmt.Sprintf("entry %d: duplicate offset %d", i, e.Offset)
			result.Errors = append(result.Errors, msg)
			result.Duplicates++
		}
		seenOffsets[e.Offset] = true

		if i > 0 {
			prev := entries[i-1]
			if !e.Timestamp.IsZero() && !prev.Timestamp.IsZero() &&
				e.Timestamp.Before(prev.Timestamp) {
				msg := fmt.Sprintf(
					"entry %d: timestamp %s is before previous %s",
					i, e.Timestamp.Format(time.RFC3339), prev.Timestamp.Format(time.RFC3339),
				)
				result.Errors = append(result.Errors, msg)
				result.OutOfOrder++
			}
		}
	}

	if !result.Valid() {
		return result, errors.New("index validation failed")
	}
	return result, nil
}
