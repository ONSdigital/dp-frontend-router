package analytics

import (
	"testing"

	jwt "github.com/form3tech-oss/jwt-go"
)

// The following code is borrowed from v3.2.2 of library:
// github.com/form3tech-oss/jwt-go
// and tweaked a little to demonstrate issues with the previously used library:
// github.com/dgrijalva/jwt-go
// BUT the code won't work with v4 of original library

func Test_mapClaims_list_aud(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": []string{"foo"},
	}
	want := true
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}
func Test_mapClaims_string_aud(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": "foo",
	}
	want := true
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func Test_mapClaims_list_aud_no_match(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": []string{"bar"},
	}
	want := false
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}
func Test_mapClaims_string_aud_fail(t *testing.T) {
	mapClaims := jwt.MapClaims{
		"aud": "bar",
	}
	want := false
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func Test_mapClaims_string_aud_no_claim(t *testing.T) {
	mapClaims := jwt.MapClaims{}
	want := false
	got := mapClaims.VerifyAudience("foo", true)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func Test_mapClaims_string_aud_no_claim_not_required(t *testing.T) {
	mapClaims := jwt.MapClaims{}
	want := false
	got := mapClaims.VerifyAudience("foo", false)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}
