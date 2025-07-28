package generator

import "crypto/rand"

func GenerateRandomSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)

	if err != nil {
		return nil, err
	}
	
	return salt, nil
}