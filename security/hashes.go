package security

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

func GenerateRandomSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash)
}

func VerifyPassword(password, hashedPassword string, salt []byte) bool {
	expectedHash := HashPassword(password, salt)
	return hashedPassword == expectedHash
}
