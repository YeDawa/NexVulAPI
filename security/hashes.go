package security

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})
}

func VerifyPassword(password, hashed string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hashed)
	return err == nil && match
}
