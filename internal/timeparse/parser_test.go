package timeparse

import (
	"testing"
	"time"
)

func TestParse_KnownFormats(t *testing.T) {
	p := New(time.UTC)

	cases := []struct {
		input    string
		wantYear int
		wantHour int
	}{
		{"2024-03-15T08:30:00Z", 2024, 8},
		{"2024-03-15T08:30:00.123456Z", 2024, 8},
		{"2024-03-15 08:30:00", 2024, 8},
		{"2024-03-15 08:30:00.000", 2024, 8},
		{"15/Mar/2024:08:30:00 +0000", 2024, 8},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := p.Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantYear)
			}
			if got.Hour() != tc.wantHour {
				t.Errorf("hour: got %d, want %d", got.Hour(), tc.wantHour)
			}
		})
	}
}

func TestParse_UnknownFormat(t *testing.T) {
	p := New(time.UTC)
	_, err := p.Parse("not-a-timestamp")
	if err == nil {
		t.Fatal("expected error for unparseable input, got nil")
	}
}

func TestAddFormat(t *testing.T) {
	p := New(time.UTC)
	p.AddFormat("01/02/2006")
	got, err := p.Parse("03/15/2024")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Month() != time.March {
		t.Errorf("month: got %v, want March", got.Month())
	}
}

func TestParseRange_Valid(t *testing.T) {
	p := New(time.UTC)
	s, e, err := p.ParseRange("2024-03-15T08:00:00Z", "2024-03-15T09:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Before(e) {
		t.Errorf("expected start before end")
	}
}

func TestParseRange_StartAfterEnd(t *testing.T) {
	p := New(time.UTC)
	_, _, err := p.ParseRange("2024-03-15T09:00:00Z", "2024-03-15T08:00:00Z")
	if err == nil {
		t.Fatal("expected error when start is after end")
	}
}
