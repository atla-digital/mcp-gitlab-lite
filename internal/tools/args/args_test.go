package args

import "testing"

func newR(m map[string]any) R { return R{m: m} }

func TestInt64_AcceptsFloat(t *testing.T) {
	r := newR(map[string]any{"pipeline_id": float64(3829)})
	if got := r.Int64("pipeline_id"); got != 3829 {
		t.Fatalf("Int64 from float64: got %d, want 3829", got)
	}
}

func TestInt64_AcceptsNumericString(t *testing.T) {
	r := newR(map[string]any{"pipeline_id": "3829"})
	if got := r.Int64("pipeline_id"); got != 3829 {
		t.Fatalf("Int64 from string: got %d, want 3829", got)
	}
}

func TestInt64_NonNumericStringReturnsZero(t *testing.T) {
	r := newR(map[string]any{"pipeline_id": "not-a-number"})
	if got := r.Int64("pipeline_id"); got != 0 {
		t.Fatalf("Int64 from garbage string: got %d, want 0", got)
	}
}

func TestInt64_MissingReturnsZero(t *testing.T) {
	r := newR(map[string]any{})
	if got := r.Int64("pipeline_id"); got != 0 {
		t.Fatalf("Int64 missing: got %d, want 0", got)
	}
}

func TestInt_AcceptsFloatAndString(t *testing.T) {
	r := newR(map[string]any{"a": float64(42), "b": "42"})
	if got := r.Int("a"); got != 42 {
		t.Fatalf("Int from float64: got %d, want 42", got)
	}
	if got := r.Int("b"); got != 42 {
		t.Fatalf("Int from string: got %d, want 42", got)
	}
}
