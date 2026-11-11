package auth

import (
	"testing"
)

func TestValidate(t *testing.T) {
	if _, e := Validate(" Test@Example.COM ", "a-long-good-password"); e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct{ email, pass string }{{"bad", "a-long-good-password"}, {"x@y.com", "short"}, {"x@y.com", ""}} {
		if _, e := Validate(tc.email, tc.pass); e == nil {
			t.Fatalf("accepted invalid credentials: %+v", tc)
		}
	}
}
func TestTokenHashStable(t *testing.T) {
	if hashToken("example") == hashToken("other") || len(hashToken("example")) != 64 {
		t.Fatal("bad token hash")
	}
}
