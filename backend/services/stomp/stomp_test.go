package stomp

import (
	"testing"

	"gopkg.in/check.v1"
)

// Runs all gocheck tests in this package.
// See other *_test.go files for gocheck tests.
func TestStomp(t *testing.T) {
	check.Suite(&StompSuite{t})
	check.TestingT(t)
}

type StompSuite struct {
	t *testing.T
}
