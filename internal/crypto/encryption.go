package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	saltLength = 16
	nonceLength = 12
	keyLength = 32
	// pbkdf2Iterations is the number of iterations for PBKDF2 key derivation
	// OWASP recommends at least 600,000 iterations for password storage
	// See: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#pbkdf2
	pbkdf2Iterations = 600_000
)

// Encrypt encrypts plaintext using AES-256-GCM with the given password
func Encrypt(plaintext, password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// Derive key from password using PBKDF2 with OWASP-recommended iterations
	key := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, keyLength, sha256.New)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Combine salt + ciphertext
	result := append(salt, ciphertext...)

	// Return base64 encoded result
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts ciphertext using AES-256-GCM with the given password
func Decrypt(ciphertext, password string) (string, error) {
	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(data) < saltLength {
		return "", errors.New("ciphertext too short")
	}

	// Extract salt and encrypted data
	salt := data[:saltLength]
	encrypted := data[saltLength:]

	// Derive key from password using PBKDF2 with OWASP-recommended iterations
	key := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, keyLength, sha256.New)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Check minimum length
	if len(encrypted) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce := encrypted[:gcm.NonceSize()]
	ciphertextData := encrypted[gcm.NonceSize():]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword checks if the provided password matches the hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
