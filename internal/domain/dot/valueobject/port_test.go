package valueobject

import "testing"

func TestPortStr(t *testing.T) {
	s := "192.168.0.1:8080"
	expected := GenerateShortURL("192_168_0_1:8080")
	if res := PortStr(s); res != expected {
		t.Errorf("portStr(%q) = %q, expected %q", s, res, expected)
	}

	s = "10.0.0.1-8080"
	expected = GenerateShortURL("10_0_0_1_8080")
	if res := PortStr(s); res != expected {
		t.Errorf("portStr(%q) = %q, expected %q", s, res, expected)
	}

	s = "foo.bar_baz"
	expected = GenerateShortURL("foo_bar_baz")
	if res := PortStr(s); res != expected {
		t.Errorf("portStr(%q) = %q, expected %q", s, res, expected)
	}
}
