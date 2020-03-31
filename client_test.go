package pushover

import "testing"

func TestClientDefaults(t *testing.T) {
	c, err := NewClient(WithAppToken("DEADBEEF"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c == nil {
		t.Fatal("expected a Client instance, got nil")
	}
	if want, have := "DEADBEEF", c.appToken; want != have {
		t.Fatalf("expected AppToken=%q, got %q", want, have)
	}
}
