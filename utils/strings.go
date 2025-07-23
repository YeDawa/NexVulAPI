package utils

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func GetLastPartOfURL(url string) string {
	parts := strings.Split(strings.TrimRight(url, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

func CountRemoteFileLines(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro ao acessar URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro: status HTTP %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("erro ao ler conteúdo: %w", err)
	}

	return lineCount, nil
}
