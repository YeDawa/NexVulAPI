package generator

import (
	"fmt"
	"time"
	"path/filepath"

	"github.com/google/uuid"
)

func Uuid(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102150405")
	
	return fmt.Sprintf("%s_%s%s", timestamp, uniqueID, ext)
}
