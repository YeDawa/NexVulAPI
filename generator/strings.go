package generator

import (
	"math/big"
	"crypto/rand"
	"encoding/hex"
	"encoding/base64"
)

func String(min, max int) string {
	length := min

	if max > min {
		randomLength, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
		length += int(randomLength.Int64())
	}

	b := make([]byte, length/2)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}

func StringBase64(n int) (string, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
