package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	FontURL  = "https://github.com/melroy89/Roboto/raw/refs/heads/main/RobotoTTF/Roboto-Regular.ttf"
	FontPath = "./temp_fonts/Roboto-Regular.ttf"
)

func DownloadFontIfNeeded() error {
	if _, err := os.Stat(FontPath); os.IsNotExist(err) {
		if err := os.MkdirAll("./temp_fonts", 0755); err != nil {
			return fmt.Errorf("failed to create font dir: %w", err)
		}

		resp, err := http.Get(FontURL)
		if err != nil {
			return fmt.Errorf("failed to download font: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("font download HTTP status: %d", resp.StatusCode)
		}

		out, err := os.Create(FontPath)
		if err != nil {
			return fmt.Errorf("failed to create font file: %w", err)
		}
		defer out.Close()

		if _, err = io.Copy(out, resp.Body); err != nil {
			return fmt.Errorf("failed to save font file: %w", err)
		}
	}
	
	return nil
}
