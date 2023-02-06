package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	return append(ciphertext, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	return origData[:(length - int(origData[length-1]))]
}

func AesEncrypt(origData, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "new cipher error", err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func AesDecrypt(encrypted string, key []byte) (text string, err error) {
	defer func() {
		if crash := recover(); crash != nil {
			err = errors.New("encrypted is aes text")
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	cryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return
	}
	origData := make([]byte, len(cryptedByte))

	blockMode.CryptBlocks(origData, cryptedByte)
	text = string(PKCS7UnPadding(origData))
	return
}
