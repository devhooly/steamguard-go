package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize       = 32
	ivSize         = 16
	keySize        = 32
	pbkdf2Iter     = 50000
)

// Encrypt encrypts data using AES-256-CBC with PBKDF2
func Encrypt(data []byte, passkey string) (encrypted []byte, ivB64 string, saltB64 string, err error) {
	// Generate salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, "", "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate IV
	iv := make([]byte, ivSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, "", "", fmt.Errorf("failed to generate IV: %w", err)
	}

	// Create key from password
	key := pbkdf2.Key([]byte(passkey), salt, pbkdf2Iter, keySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Add padding
	paddedData := pkcs7Pad(data, aes.BlockSize)

	// Encrypt
	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	// Encode to base64
	ivB64 = base64.StdEncoding.EncodeToString(iv)
	saltB64 = base64.StdEncoding.EncodeToString(salt)

	return ciphertext, ivB64, saltB64, nil
}

// Decrypt decrypts data using AES-256-CBC with PBKDF2
func Decrypt(encrypted []byte, passkey string, ivB64 string, saltB64 string) ([]byte, error) {
	// Decode IV and salt from base64
	iv, err := base64.StdEncoding.DecodeString(ivB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Create key from password
	key := pbkdf2.Key([]byte(passkey), salt, pbkdf2Iter, keySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Check size
	if len(encrypted)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid encrypted data size")
	}

	// Decrypt
	plaintext := make([]byte, len(encrypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, encrypted)

	// Remove padding
	unpaddedData, err := pkcs7Unpad(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to remove padding: %w", err)
	}

	return unpaddedData, nil
}

// pkcs7Pad adds PKCS7 padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad removes PKCS7 padding
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding")
	}

	// Verify padding correctness
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}

// SecureZero wipes data in memory
func SecureZero(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
