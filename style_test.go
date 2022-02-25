package browser

import "testing"

func TestExample(t *testing.T) {
	s := Size{Value: 100, Unit: UnitPX}

	if got, want := s.String(), "100px"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
