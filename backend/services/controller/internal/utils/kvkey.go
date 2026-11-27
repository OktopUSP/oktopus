package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// EncodeKVKey converts a credential username (e.g. a USP Endpoint ID such as
// "os::000F94-AM.123456") into a valid NATS KeyValue key.
//
// NATS KV keys only accept [-/_=.a-zA-Z0-9] and cannot start/end with '.' or
// contain "..". Every other byte, plus '=' itself, is escaped as "=HH" (upper
// case hex), so the encoding is reversible and plain usernames made only of
// letters, digits, '-', '_', '/' and inner '.' are kept unchanged.
func EncodeKVKey(name string) string {
	var b strings.Builder
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9',
			c == '-', c == '_', c == '/':
			b.WriteByte(c)
		case c == '.' && i > 0 && i < len(name)-1 && name[i-1] != '.' && name[i+1] != '.':
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "=%02X", c)
		}
	}
	return b.String()
}

// DecodeKVKey reverses EncodeKVKey. Keys that are not valid encodings are
// returned unchanged.
func DecodeKVKey(key string) string {
	var b strings.Builder
	for i := 0; i < len(key); i++ {
		if key[i] != '=' {
			b.WriteByte(key[i])
			continue
		}
		if i+2 >= len(key) {
			return key
		}
		v, err := strconv.ParseUint(key[i+1:i+3], 16, 8)
		if err != nil {
			return key
		}
		b.WriteByte(byte(v))
		i += 2
	}
	return b.String()
}
