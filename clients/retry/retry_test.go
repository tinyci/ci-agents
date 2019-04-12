package retry

import (
	"testing"

	"github.com/erikh/check"
)

type retrySuite struct{}

var _ = check.Suite(&retrySuite{})

func TestRetry(t *testing.T) {
	check.TestingT(t)
}
