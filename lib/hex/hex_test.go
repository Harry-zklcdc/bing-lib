package hex_test

import (
	"testing"

	"github.com/Harry-zklcdc/bing-lib/lib/hex"
)

func TestHex(t *testing.T) {
	t.Log(hex.NewHex(32))
	t.Log(hex.NewHexLowercase(32))
}
