package hex_test

import (
	"testing"

	"github.com/Harry-zklcdc/bing-lib/lib/hex"
)

func TestUUID(t *testing.T) {
	t.Log(hex.NewUUID())
}
