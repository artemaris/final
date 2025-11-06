package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "This is a test message"
	password := "test-password-123"

	// Encrypt
	ciphertext, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	if ciphertext == "" {
		t.Fatal("Ciphertext is empty")
	}

	// Decrypt
	decrypted, err := Decrypt(ciphertext, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original. Got: %s, Want: %s", decrypted, plaintext)
	}
}

func TestEncryptDecryptWrongPassword(t *testing.T) {
	plaintext := "Secret message"
	password := "correct-password"
	wrongPassword := "wrong-password"

	ciphertext, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to decrypt with wrong password
	_, err = Decrypt(ciphertext, wrongPassword)
	if err == nil {
		t.Fatal("Expected decryption with wrong password to fail, but it succeeded")
	}
}

func TestHashPassword(t *testing.T) {
	password := "my-secure-password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}

	if hash == "" {
		t.Fatal("Hash is empty")
	}

	// Same password should produce different hash
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}

	if hash == hash2 {
		t.Fatal("Same password produced identical hash (should be different due to salt)")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "my-password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}

	// Correct password
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword failed for correct password")
	}

	// Wrong password
	if CheckPassword("wrong-password", hash) {
		t.Error("CheckPassword succeeded for wrong password")
	}
}

func TestEncryptEmptyString(t *testing.T) {
	ciphertext, err := Encrypt("", "password")
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, "password")
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != "" {
		t.Errorf("Expected empty string, got: %s", decrypted)
	}
}
