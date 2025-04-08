package utility

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var key = []byte("ThIsis32bYteKeyForAES256exAmple!")

func pad(data []byte) []byte {
	padLen := aes.BlockSize - len(data)%aes.BlockSize
	return append(data, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
}

func unpad(data []byte) []byte {
	return data[:len(data)-int(data[len(data)-1])]
}

func Encrypt(text string) (string, error) {
	plain := pad([]byte(text))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	_, _ = io.ReadFull(rand.Reader, iv)

	cipherText := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText, plain)

	final := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(final), nil
}

func Decrypt(encrypted string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(encrypted)
	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plain := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, cipherText)

	return string(unpad(plain)), nil
}
