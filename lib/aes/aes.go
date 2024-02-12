package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

// 填充明文
func pkcs7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// 去除填充数据
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES 加密
func Encrypt(Msg string, Key string) (string, error) {
	origData := []byte(Msg)
	key := []byte(Key)
	if len(key) != 32 {
		return "", errors.New("key length must be 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	origData = pkcs7Padding(origData, blockSize)
	cipherText := make([]byte, aes.BlockSize+len(origData))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherText[aes.BlockSize:], origData)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 解密
func Decrypt(Msg string, Key string) (string, error) {
	Msg = strings.ReplaceAll(Msg, " ", "+")
	cipherText, _ := base64.StdEncoding.DecodeString(Msg)
	key := []byte(Key)
	if len(key) != 32 {
		return "", errors.New("key length must be 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = pkcs7UnPadding(origData)
	return string(origData), nil
}

// func Encrypt(e, t string) ([]byte, error) {
// 	block, err := aes.NewCipher([]byte(t))
// 	if err != nil {
// 		return []byte{}, err
// 	}

// 	ciphertext := make([]byte, aes.BlockSize+len(e))
// 	iv := ciphertext[:aes.BlockSize]
// 	mode := cipher.NewCBCEncrypter(block, iv)
// 	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(e))

// 	return ciphertext, nil
// }

// func Decrypt(e, t string) (string, error) {
// 	block, err := aes.NewCipher([]byte(t))
// 	if err != nil {
// 		return "", err
// 	}

// 	ciphertext, _ := base64.StdEncoding.DecodeString(e)
// 	if len(ciphertext) < aes.BlockSize {
// 		return "", fmt.Errorf("ciphertext too short")
// 	}
// 	iv := ciphertext[:aes.BlockSize]
// 	ciphertext = ciphertext[aes.BlockSize:]

// 	mode := cipher.NewCBCDecrypter(block, iv)
// 	mode.CryptBlocks(ciphertext, ciphertext)

// 	return string(ciphertext), nil
// }
