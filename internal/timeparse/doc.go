// Package timeparse provides flexible timestamp parsing for log lines.
//
// It supports a wide range of common log timestamp formats and allows
// callers to register custom formats. The Parser type is the primary
// entry point; create one with New and optionally add extra formats
// with AddFormat before calling Parse or ParseRange.
//
// Example:
//
//	p := timeparse.New(time.UTC)
//	t, err := p.Parse("2024-03-15T08:30:00Z")
//
//	start, end, err := p.ParseRange(
//		"2024-03-15T08:00:00Z",
//		"2024-03-15T09:00:00Z",
//	)
package timeparse
