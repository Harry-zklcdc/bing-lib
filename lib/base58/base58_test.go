package base58_test

import (
	"testing"

	"github.com/Harry-zklcdc/bing-lib/lib/base58"
)

const STR = "Harry-zklcdc/go-proxy-bingai"

func TestBase58(t *testing.T) {
	base58Str := base58.Base58Encoding(STR)
	t.Log("Base58Encode", base58Str)

	originStr := base58.Base58Decoding(base58Str)
	t.Log("Base58Decode", originStr)

	if originStr != STR {
		t.Error("Base58Decode failed")
	}
}
