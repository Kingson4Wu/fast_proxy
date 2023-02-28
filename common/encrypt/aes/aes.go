package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

/**
https://www.jianshu.com/p/0caab60fea9f
*/
// Encryption process:
// 1. Process the data and pad it using PKCS7 (when the key length is insufficient, fill in the corresponding number of bytes).
// 2. Encrypt the data using AES encryption with CBC mode.
// 3. Encrypt the resulting encrypted data using base64 and obtain a string.
// The decryption process is the opposite.

// If the string is 16, 24, or 32 characters long, it corresponds to the AES-128, AES-192, and AES-256 encryption methods, respectively.
// The key must not be leaked.

// pkcs7Padding padding
func pkcs7Padding(data []byte, blockSize int) []byte {
	// Determine the number of missing bytes. At least 1 and at most blockSize
	padding := blockSize - len(data)%blockSize
	// Fill in the missing bytes. Copy the slice []byte{byte(padding)} padding times
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding is the reverse operation of padding
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("encrypted string error")
	}
	// Get the number of padding bytes
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	// Create an encryption instance
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes encrypt error: %w", err)
	}
	// Determine the block size for encryption
	blockSize := block.BlockSize()
	// Padding
	encryptBytes := pkcs7Padding(data, blockSize)
	// Initialize the slice to receive the encrypted data
	encrypted := make([]byte, len(encryptBytes))
	// Use CBC mode
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// Execute encryption
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return encrypted, nil
}

func Decrypt(data []byte, key []byte) ([]byte, error) {
	// Create an instance
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes decrypt error: %w", err)
	}
	// Get the block size
	blockSize := block.BlockSize()
	// Use CBC
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// Initialize the slice to receive the decrypted data
	decrypted := make([]byte, len(data))
	// Execute decryption
	blockMode.CryptBlocks(decrypted, data)
	// Remove padding
	decrypted, err = pkcs7UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

func EncryptByAesWithKey(data string, pwd string) (string, error) {
	res, err := Encrypt([]byte(data), []byte(pwd))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func DecryptByAesWithKey(data string, pwd string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	result, err := Decrypt(dataByte, []byte(pwd))
	if err != nil {
		return "", err
	}
	return string(result), nil
}
