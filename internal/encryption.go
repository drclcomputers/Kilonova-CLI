// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

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

var masterKey = []byte("ThIsis32bYteKeyForAES256exAmple!")

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
	plainText := []byte(text)

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", fmt.Errorf("cipher initialization failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM mode failed: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce generation failed: %w", err)
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", fmt.Errorf("cipher initialization failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM mode failed: %w", err)
	}

	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherText := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed (integrity check failed): %w", err)
	}

	return string(plainText), nil
}
