package utils

import (
	"fmt"
	"strings"
)

func StringTemplateReplace(s string, templates map[string]string) string {
	newS := s

	for key, val := range templates {
		t := fmt.Sprintf("{{%s}}", key)
		newS = strings.ReplaceAll(newS, t, val)
	}

	return newS
}

// func processFlagMap() {}
