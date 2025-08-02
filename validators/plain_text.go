package validators

import (
	"fmt"
	"net/http"
	"strings"
)

func ValidateTextPlainURL(url string) error {
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		return fmt.Errorf("invalid content type: %s", contentType)
	}

	return nil
}
