package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// EncryptAndHash encrypts plaintext (≤配置的 50MB 级别) using AES-GCM with a fresh nonce.
// It returns a reader over nonce+ciphertext, exact size, and SHA-256 hex of the full ciphertext blob.
type cipherResult struct {
	reader io.Reader
	size   int64
	hash   string
}

func (c cipherResult) Reader() io.Reader { return c.reader }
func (c cipherResult) Size() int64       { return c.size }
func (c cipherResult) Hash() string      { return c.hash }

func EncryptAndHash(r io.Reader, key []byte) (cipherResult, error) {
	if len(key) != 32 {
		return cipherResult{}, fmt.Errorf("key length must be 32 bytes for AES-256-GCM")
	}
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return cipherResult{}, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return cipherResult{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return cipherResult{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return cipherResult{}, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	blob := append(nonce, ciphertext...)
	hash := sha256.Sum256(blob)
	return cipherResult{
		reader: bytes.NewReader(blob),
		size:   int64(len(blob)),
		hash:   hex.EncodeToString(hash[:]),
	}, nil
}

// DecryptReader decrypts AES-GCM ciphertext (nonce prefixed) fully in-memory and returns a reader for plaintext.
type decryptResult struct {
	reader io.Reader
}

func (d decryptResult) Reader() io.Reader { return d.reader }

func DecryptReader(r io.Reader, key []byte) (decryptResult, error) {
	if len(key) != 32 {
		return decryptResult{}, fmt.Errorf("key length must be 32 bytes for AES-256-GCM")
	}
	blob, err := io.ReadAll(r)
	if err != nil {
		return decryptResult{}, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return decryptResult{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return decryptResult{}, err
	}
	nonceSize := gcm.NonceSize()
	if len(blob) < nonceSize {
		return decryptResult{}, fmt.Errorf("ciphertext too short")
	}
	nonce := blob[:nonceSize]
	ciphertext := blob[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return decryptResult{}, err
	}
	return decryptResult{reader: bytes.NewReader(plaintext)}, nil
}
