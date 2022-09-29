package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// GetTemplateDir get template path
// return err if template file doesn't exist
func GetTemplateDir(template string) (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")
	templatePath := filepath.Join(basePath, fmt.Sprintf("sangtool/templates/%s.json", template))
	// check template exist
	_, err := os.Stat(templatePath)
	if os.IsNotExist(err) {
		return "", errors.New(fmt.Sprintf("template:%s does not exist", template))
	}
	return templatePath, nil
}
