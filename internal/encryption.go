package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var key = []byte("ThIsis32bYteKeyForAES256exAmple!")

func pad(data []byte) []byte {
	padLen := aes.BlockSize - len(data)%aes.BlockSize
	return append(data, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
}

func unpad(data []byte) ([]byte, error) {
	paddingLen := int(data[len(data)-1])
	if paddingLen > aes.BlockSize || paddingLen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-paddingLen], nil
}

func Encrypt(text string) (string, error) {
	plain := pad([]byte(text))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(plain))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(cipherText, plain)

	final := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(final), nil
}

func Decrypt(encrypted string) (string, error) {

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plain := make([]byte, len(cipherText))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plain, cipherText)

	plain, err = unpad(plain)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
