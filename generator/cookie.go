package generator

import (
	"fmt"
	"encoding/base64"
)

func CookieValue(userID uint) (string, error) {
	randomString, err := StringBase64(32)

	if err != nil {
		return "", err
	}

	cookieValue := fmt.Sprintf("%s:%d", randomString, userID)
	return base64.StdEncoding.EncodeToString([]byte(cookieValue)), nil
}
