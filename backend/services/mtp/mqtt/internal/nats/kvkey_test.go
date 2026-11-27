package nats

import "testing"

func TestEncodeKVKey(t *testing.T) {
	tests := map[string]string{
		"oktopusController":    "oktopusController",
		"00256D-0123456789":    "00256D-0123456789",
		"my.device_1/a-b":      "my.device_1/a-b",
		"os::000F94-AM.123456": "os=3A=3A000F94-AM.123456",
		"user@example.com":     "user=40example.com",
		"a=b":                  "a=3Db",
		".leading":             "=2Eleading",
		"trailing.":            "trailing=2E",
		"double..dot":          "double=2E=2Edot",
	}
	for in, want := range tests {
		got := EncodeKVKey(in)
		if got != want {
			t.Errorf("EncodeKVKey(%q) = %q, want %q", in, got, want)
		}
	}
}
