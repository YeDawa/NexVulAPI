package utils

import (
	"bufio"
	"fmt"
	"net/http"
)

func StreamRemoteFile(url string, lineHandler func(string)) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to perform GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		lineHandler(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	return nil
}
