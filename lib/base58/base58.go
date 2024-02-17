package base58

import (
	"bytes"
	"math/big"
)

var base58 = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encoding(str string) string {
	strByte := []byte(str)
	strTen := big.NewInt(0).SetBytes(strByte)

	var modSlice []byte
	for strTen.Cmp(big.NewInt(0)) > 0 {
		mod := big.NewInt(0)
		strTen58 := big.NewInt(58)
		strTen.DivMod(strTen, strTen58, mod)
		modSlice = append(modSlice, base58[mod.Int64()])
	}

	for _, elem := range strByte {
		if elem != 0 {
			break
		} else if elem == 0 {
			modSlice = append(modSlice, byte('1'))
		}
	}

	ReverseModSlice := reverseByteArr(modSlice)
	return string(ReverseModSlice)
}

func Base58Decoding(str string) string {
	strByte := []byte(str)
	ret := big.NewInt(0)
	for _, byteElem := range strByte {
		index := bytes.IndexByte(base58, byteElem)
		ret.Mul(ret, big.NewInt(58))
		ret.Add(ret, big.NewInt(int64(index)))
	}

	return string(ret.Bytes())
}

func reverseByteArr(bytes []byte) []byte { //将字节的数组反转
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i] //前后交换
	}
	return bytes
}
