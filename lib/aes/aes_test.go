package aes_test

import (
	"testing"

	"github.com/Harry-zklcdc/bing-lib/lib/aes"
)

func TestAES(t *testing.T) {
	t.Log("TestEncrypt")
	c, err := aes.Encrypt("Harry-zklcdc/go-proxy-bingai", "NFKP6NEN0TH50YD1HMKVBW5EOVGFYMBL")
	if err != nil {
		t.Error(err)
	}
	t.Log(c)

	t.Log("TestDecrypt")
	d, err := aes.Decrypt(c, "NFKP6NEN0TH50YD1HMKVBW5EOVGFYMBL")
	if err != nil {
		t.Error(err)
	}
	t.Log(d)
}
