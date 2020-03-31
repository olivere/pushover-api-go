package pushover

import "testing"

func TestCut(t *testing.T) {
	tests := []struct {
		Input    string
		Max      int
		Trailing string
		Want     string
	}{
		{"Oliver", 10, "", "Oliver"},
		{"Oliver", 5, ellipsis, "Oliv…"},
		{"Olli", 5, ellipsis, "Olli"},
		{"", 5, ellipsis, ""},
	}
	for _, tt := range tests {
		if want, have := tt.Want, cut(tt.Input, tt.Max, tt.Trailing); want != have {
			t.Fatalf("want cut(%q, %d, %q)=%q, have %q", tt.Input, tt.Max, tt.Trailing, want, have)
		}
	}
}
