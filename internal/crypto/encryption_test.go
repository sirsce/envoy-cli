package crypto

import (
	"bytes"
	"testing"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{"AES-128", 16, false},
		{"AES-192", 24, false},
		{"AES-256", 32, false},
		{"Invalid size", 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)
			_, err := NewEncryptor(key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"Simple text", []byte("DATABASE_URL=postgres://localhost:5432/db")},
		{"Multi-line env", []byte("API_KEY=secret123\nDB_PASS=pass456\n")},
		{"Empty string", []byte("")},
		{"Special chars", []byte("KEY=!@#$%^&*()_+-={}[]|:;<>?,./")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := encryptor.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			decrypted, err := encryptor.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Decrypted text doesn't match original.\nGot: %s\nWant: %s", decrypted, tt.plaintext)
			}
		})
	}
}

func TestDecryptInvalidData(t *testing.T) {
	key, _ := GenerateKey(32)
	encryptor, _ := NewEncryptor(key)

	tests := []struct {
		name       string
		ciphertext string
	}{
		{"Invalid base64", "not-valid-base64!@#"},
		{"Too short", "YWJj"},
		{"Wrong encryption", "dGVzdGRhdGE="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tt.ciphertext)
			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	key1, err := GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	key2, err := GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	if bytes.Equal(key1, key2) {
		t.Error("Generated keys should be unique")
	}
}
