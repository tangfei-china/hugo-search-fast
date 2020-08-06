package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

var key = []byte("7qaz8wsx")


func DESCBCEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	cipherText := make([]byte, len(origData))
	blockMode.CryptBlocks(cipherText, origData)
	return cipherText, nil
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func DESCBCDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := cipherText
	blockMode.CryptBlocks(origData, cipherText)

	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//加密内容
func DESEncryptString(str string) (string, error) {
	input := []byte(str)
	encrypt, e := DESCBCEncrypt(input, key)
	if e != nil {
		return "", e
	}

	return base64.StdEncoding.EncodeToString(encrypt), nil
}

//解密内容
func DESDecryptString(str string) (string, error) {
	input, e := base64.StdEncoding.DecodeString(str)
	if e != nil {
		return "", e
	}
	decrypt, e := DESCBCDecrypt(input, key)
	if e != nil {
		return "", e
	}
	return string(decrypt), nil
}
