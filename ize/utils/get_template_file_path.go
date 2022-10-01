package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetTemplateFilePath(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	toolDir := filepath.Join(home, "ize_templates")
	os.Mkdir(toolDir, 0777)

	exPath := filepath.Join(toolDir, fmt.Sprintf("%s.json", name))

	return exPath, nil
}
