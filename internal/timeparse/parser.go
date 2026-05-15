package timeparse

import (
	"fmt"
	"time"
)

// Common log timestamp formats to attempt parsing.
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05.000000",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.000000",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// Parser holds configuration for timestamp parsing.
type Parser struct {
	formats  []string
	location *time.Location
}

// New creates a Parser with default formats and the given timezone.
// If loc is nil, time.UTC is used.
func New(loc *time.Location) *Parser {
	if loc == nil {
		loc = time.UTC
	}
	return &Parser{
		formats:  knownFormats,
		location: loc,
	}
}

// AddFormat registers an additional custom timestamp format.
func (p *Parser) AddFormat(format string) {
	p.formats = append(p.formats, format)
}

// Parse attempts to parse s using all known formats.
// Returns the first successful result or an error if none match.
func (p *Parser) Parse(s string) (time.Time, error) {
	for _, f := range p.formats {
		if t, err := time.ParseInLocation(f, s, p.location); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timeparse: cannot parse %q with any known format", s)
}

// ParseRange parses a start and end timestamp string.
// Returns an error if either cannot be parsed or if start is after end.
func (p *Parser) ParseRange(start, end string) (time.Time, time.Time, error) {
	s, err := p.Parse(start)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("start: %w", err)
	}
	e, err := p.Parse(end)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("end: %w", err)
	}
	if s.After(e) {
		return time.Time{}, time.Time{}, fmt.Errorf("timeparse: start %v is after end %v", s, e)
	}
	return s, e, nil
}
