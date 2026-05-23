package index

import (
	"testing"
	"time"
)

func makeValidEntries() []Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: base, Offset: 0},
		{Timestamp: base.Add(time.Minute), Offset: 100},
		{Timestamp: base.Add(2 * time.Minute), Offset: 200},
	}
}

func TestValidate_ValidEntries(t *testing.T) {
	entries := makeValidEntries()
	result, err := Validate(entries)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !result.Valid() {
		t.Errorf("expected valid result, got errors: %v", result.Errors)
	}
	if result.TotalEntries != 3 {
		t.Errorf("expected 3 entries, got %d", result.TotalEntries)
	}
}

func TestValidate_EmptyEntries(t *testing.T) {
	result, err := Validate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalEntries != 0 {
		t.Errorf("expected 0 entries")
	}
}

func TestValidate_ZeroTimestamp(t *testing.T) {
	entries := makeValidEntries()
	entries[1].Timestamp = time.Time{}
	result, err := Validate(entries)
	if err == nil {
		t.Fatal("expected error for zero timestamp")
	}
	if result.InvalidCount != 1 {
		t.Errorf("expected 1 invalid, got %d", result.InvalidCount)
	}
}

func TestValidate_OutOfOrder(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: base.Add(time.Minute), Offset: 0},
		{Timestamp: base, Offset: 100},
	}
	result, err := Validate(entries)
	if err == nil {
		t.Fatal("expected error for out-of-order entries")
	}
	if result.OutOfOrder != 1 {
		t.Errorf("expected 1 out-of-order, got %d", result.OutOfOrder)
	}
}

func TestValidate_DuplicateOffset(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: base, Offset: 0},
		{Timestamp: base.Add(time.Minute), Offset: 0},
	}
	result, err := Validate(entries)
	if err == nil {
		t.Fatal("expected error for duplicate offset")
	}
	if result.Duplicates != 1 {
		t.Errorf("expected 1 duplicate, got %d", result.Duplicates)
	}
}

func TestValidate_NegativeOffset(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: base, Offset: -1},
	}
	result, err := Validate(entries)
	if err == nil {
		t.Fatal("expected error for negative offset")
	}
	if result.InvalidCount != 1 {
		t.Errorf("expected 1 invalid, got %d", result.InvalidCount)
	}
}
